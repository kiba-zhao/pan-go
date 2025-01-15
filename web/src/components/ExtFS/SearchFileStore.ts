
import type { Store,StoreContext } from "./Store";
import { newStore,initStoreContext ,emitChange} from "./Store";

import type {API,ExtFSSearchFile,ExtFSSearchFileSearchCondition,ExtFSRemoteSearchFileSearchCondition} from "../../api"

type SearchFileItem = { peerId?: string } & ExtFSSearchFile;


type SearchFileStoreData = {
    files:SearchFileItem[]
    isComplete:boolean
}

type SearchFileContext =  {
    isSyncNodeComplete:boolean
    isSyncRemoteComplete:boolean
    abortCtrl:AbortController
    workers:Promise<void>[]
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
        ctx.data = data;
    }
    ctx.workers = [
        syncFromNode(ctx,query,api),
        syncFromRemotes(ctx,query,api)
    ];

}

async function syncFromNode(ctx:SearchFileContext,query:string,api:API){
    const {data} = ctx;
    const generator = generateFromNode(query,api);
    for await (const files of generator) {
        if (ctx.data !== data)
            break
        data.files = [...data.files,...files];
        emitChange(ctx);
    }
    if (ctx.data !== data) return
    ctx.isSyncNodeComplete = true;
    ctx.data.isComplete = ctx.isSyncNodeComplete && ctx.isSyncRemoteComplete;
    if (ctx.data.isComplete) {
        emitChange(ctx);
    }
}

async function *generateFromNode(query:string,api:API):AsyncGenerator<ExtFSSearchFile[]> {
    const condition = {} as ExtFSSearchFileSearchCondition
    condition.query = query
    condition._start = 0
    condition._end = 100

    while(true) {
        const [hash,total,files] = await api.searchExtFSSearchFileResults(condition);
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
    }
    
}

async function syncFromRemotes(ctx:SearchFileContext,query:string,api:API){
    const {data} = ctx;
    const remotes = await api.selectAllExtFSRemoteNodes();
    if (remotes.length > 0) {
        const generators = {} as Record<string,AsyncGenerator<ExtFSSearchFile[]> | null>;
        let offset = 0;
        let pendingNum = remotes.length;
        while (true){
            // init generator by peerId
            const {peerId} = remotes[offset];
            offset = offset+1 % remotes.length;
            if (generators[peerId] === null){
                continue;
            }
            if (generators[peerId] === void 0){
                generators[peerId] = generateFromRemote(peerId,query,api);
            }
            const generator = generators[peerId];
            // 

            // generate SearchFileItem
            const {value,done} = await generator.next()
            if (ctx.data !== data)
                break
            const files = (value as ExtFSSearchFile[]).map(_=>({..._,peerId}))
            data.files = [...data.files,...files];
            emitChange(ctx);
            // 

            if (done){
                generators[peerId] = null
                pendingNum--;
            }
            if (pendingNum <=0){
                break
            }
        }
    }

    if (ctx.data !== data) return
    ctx.isSyncRemoteComplete = true;
    ctx.data.isComplete = ctx.isSyncNodeComplete && ctx.isSyncRemoteComplete;
    if (ctx.data.isComplete) {
        emitChange(ctx);
    }
}

async function *generateFromRemote(peerId:string,query:string,api:API):AsyncGenerator<ExtFSSearchFile[]> {
    const condition = {} as ExtFSRemoteSearchFileSearchCondition
    condition.peerId = peerId
    condition.query = query
    condition._start = 0
    condition._end = 100

    while(true) {
        const [hash,total,files] = await api.searchExtFSRemoteSearchFileResults(condition);
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
    }
}