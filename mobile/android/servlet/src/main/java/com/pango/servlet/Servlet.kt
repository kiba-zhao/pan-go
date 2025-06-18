package com.pango.servlet

import android.content.Context
import com.facebook.react.ReactPackage

class Servlet {

    private val mAgentProxy = AgentProxy()
    val reactPackage: ReactPackage = ReactNativePackage(mAgentProxy)

    fun init(ctx: Context){
        mAgentProxy.init(ctx)
    }

    fun start(){
        mAgentProxy.run()
    }

    fun stop(){
        mAgentProxy.terminate()
    }

    companion object {
        fun checkPermissions(context: Context):Pair<Boolean,Array<String>>{
            return AgentProxy.checkPermissions(context)
        }
    }
}