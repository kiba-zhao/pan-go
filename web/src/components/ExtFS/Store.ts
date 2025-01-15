
export type StoreListener = ()=>void
export interface Store <T extends any> {
    subscribe(listener:StoreListener):StoreListener
    getSnapshot():T
}

export type StoreContext<T extends any> = {
    listeners:StoreListener[]
    data:T
}

export function newStore<T extends any>(ctx:StoreContext<T>):Store<T>{
    return {
        subscribe:(listener)=>subscribe(ctx,listener),
        getSnapshot:()=>getSnapshot(ctx)
    }
}

export function subscribe<T extends any>(ctx : StoreContext<T>,listener:StoreListener):StoreListener{
    ctx.listeners = [...ctx.listeners, listener];
    return () => {
        ctx.listeners = ctx.listeners.filter(l => l !== listener);
    }
}

export function getSnapshot<T extends any>(ctx : StoreContext<T>):T{
    return ctx.data;
}

export function emitChange<T extends any>(ctx : StoreContext<T>) {
    for (const listener of ctx.listeners) {
        listener();
    }
}


export function initStoreContext<T extends any>(ctx : StoreContext<T>,data:T){
    ctx.data = data;
    ctx.listeners = [];
}