package com.pango.servlet

import android.util.Log
import com.facebook.react.bridge.Promise
import com.facebook.react.bridge.ReactApplicationContext
import com.facebook.react.bridge.ReactContextBaseJavaModule
import com.facebook.react.bridge.ReactMethod


internal class AgentNativeModule (reactContext: ReactApplicationContext,agentProxy: AgentProxy) :
    ReactContextBaseJavaModule(reactContext) {
    private val mAgentProxy = agentProxy

    @ReactMethod
    fun execWithJSON(action:String,body:String,promise:Promise){
        Log.d(NAME, "execWithJSON begin: $action $body")
        
        try {
            val results = mAgentProxy.exec(action.toByteArray(),body.toByteArray())
            promise.resolve(String(results))
        }catch (e:Exception){
            Log.e(NAME, "execWithJSON Error: ${e.toString()}")
            promise.reject(e)
        }
        Log.d(NAME, "execWithJSON end: $action $body")
    }


    override fun getName() = NAME

    companion object {
        const val NAME = "AgentNativeModule"
    }
}