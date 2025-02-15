
import type { Store,StoreContext } from "./Store";
import { newStore,initStoreContext ,emitChange} from "./Store";

import type {API,ExtFSSearchFile,ExtFSSearchFileSearchCondition,ExtFSRemoteSearchFileSearchCondition, ExtFSRemoteNode} from "../../api"

type SearchFileItem = { peerId?: string } & ExtFSSearchFile;


type SearchFileStoreData = {
    files:SearchFileItem[]
    isComplete:boolean
    errs:Record<string,any>
}

type SearchFileContext =  {
    workerId:Symbol
    isSyncNodeComplete:boolean
    isSyncRemoteComplete:boolean
    abortCtrl:AbortController
    worker:Promise<void>
}& StoreContext<SearchFileStoreData>

export interface SearchFileStore extends Store<SearchFileStoreData> {
    abort(reason?: any):void
    refresh():void
}

export function newSearchFileStore(query:string,api:API): SearchFileStore {
    const ctx = {} as SearchFileContext
    const data = {} as SearchFileStoreData
    data.files = [];
    initStoreContext(ctx,data);
    refresh(ctx,query,api);

    const store = newStore(ctx);
    return {
        ...store,
        abort :(reason?:any)=>abort(ctx,reason),
        refresh:()=>refresh(ctx,query,api)
    }
}

function abort(ctx:SearchFileContext,reason?: any){
    ctx.abortCtrl && ctx.abortCtrl.abort(reason)
}

function refresh(ctx:SearchFileContext,query:string,api:API){
    ctx.abortCtrl = new AbortController()
    ctx.isSyncNodeComplete = false;
    ctx.isSyncRemoteComplete = false;
    if (ctx.data.files.length>0 || ctx.data.isComplete){
        const data = {} as SearchFileStoreData
        data.files = [];
        data.errs = {};
        ctx.data = data;
    }
    ctx.workerId = Symbol();
    ctx.worker = sync(ctx.workerId,ctx,query,api);

}

async function sync(workerId:Symbol,ctx:SearchFileContext,query:string,api:API){
    if (workerId !== ctx.workerId) return;
    const generator = generateSearchFiles(query,api,ctx.abortCtrl.signal);
    const workers =  [flushWithGenerator(workerId,ctx,generator)];

    const remotes = await api.selectAllExtFSRemoteNodes();
    if (remotes.length > 0) {
        workers.push(syncRemotes(workerId,ctx,query,api,remotes));
    }

    await Promise.all(workers);

    if (workerId !== ctx.workerId) return;
    ctx.data = {...ctx.data,isComplete:true};
    emitChange(ctx);
}

async function syncRemotes(workerId:Symbol,ctx:SearchFileContext,query:string,api:API,remotes:ExtFSRemoteNode[],offset:number =0){
    if (workerId !== ctx.workerId) return;
    const {peerId} = remotes[offset];
    const generator = generateSearchFiles(query,api,ctx.abortCtrl.signal,peerId);
    await flushWithGenerator(workerId,ctx,generator,peerId);
    if (offset < remotes.length - 1){
        await syncRemotes(workerId,ctx,query,api,remotes,offset+1);
    }
}

async function flushWithGenerator(workerId:Symbol,ctx:SearchFileContext,generator:AsyncGenerator<ExtFSSearchFile[]>,peerId?:string){
    if (workerId !== ctx.workerId) return;
    while (true){
        const {value,done} = await generator.next()
        if (workerId !== ctx.workerId) break;
        if (done){
            if (value){
                ctx.data.errs = {...ctx.data.errs,[peerId||""]:value};
                emitChange(ctx);
            }
            break;
        }
        let files = value as ExtFSSearchFile[]
        if (files.length <= 0){
            continue;
        }
        if (peerId){
            files = files.map(_=>({..._,peerId}))
        }
        let files_ = [...ctx.data.files,...files];
        ctx.data = {...ctx.data,files:files_};
        emitChange(ctx);
    }
}

async function *generateSearchFiles(query:string,api:API,signal:AbortSignal,peerId?:string):AsyncGenerator<ExtFSSearchFile[]> {
    
    const condition = {} as ExtFSSearchFileSearchCondition
    condition.query = query
    condition._start = 0
    condition._end = 100
    let err;

    let hash:string;
    let total:number;
    let files:ExtFSSearchFile[];
    while(true) {
        try{
            if (peerId) {
                [hash,total,files] = await api.searchExtFSRemoteSearchFileResults(peerId,condition as ExtFSRemoteSearchFileSearchCondition,{signal});

            }else{
                [hash,total,files] = await api.searchExtFSSearchFileResults(condition,{signal});
            }
       
            yield files;
            if (total ===0) {
                break;
            }
        
            if (condition.hash === void 0) {
                condition.hash = hash;
            }else if (condition.hash !== hash) {
                throw new Error("hash not match");
            }

            condition._start = condition._start + files.length;
            if (total >0 && condition._start>= total) {
                break;
            }
            condition._end = condition._end + files.length;
            if (total >0 && condition._end > total){
                condition._end = total;
            }

            if (total <0){
                await (new Promise(resolve => setTimeout(resolve, 1500)))
            }
        }catch(e){
            err = e;
            break;
        }
    }
    
    return err;
}