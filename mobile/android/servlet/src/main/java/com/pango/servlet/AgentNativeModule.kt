package com.pango.servlet

import android.util.Log
import com.alibaba.fastjson2.JSON
import com.facebook.react.bridge.Arguments
import com.facebook.react.bridge.Promise
import com.facebook.react.bridge.ReactApplicationContext
import com.facebook.react.bridge.ReactContextBaseJavaModule
import com.facebook.react.bridge.ReactMethod
import com.facebook.react.bridge.ReadableMap
import gomobile.GoMobileAgent


class AgentNativeModule (reactContext: ReactApplicationContext,agent: GoMobileAgent?) :
    ReactContextBaseJavaModule(reactContext) {
    private val mAgent:GoMobileAgent? = agent

    @ReactMethod
    fun execWithJSON(action:String,body:String,promise:Promise){
        Log.d(NAME, "execWithJSON begin: $action $body")
        if (mAgent == null) {
            promise.reject(Exception("Not Available"))
            return
        }

        try {
            val results = mAgent.exec(action.toByteArray(),body.toByteArray())
            promise.resolve(String(results))
        }catch (e:Exception){
            Log.e(NAME, "execWithJSON Error: ${e.toString()}")
            promise.reject(e)
        }
        Log.d(NAME, "execWithJSON end: $action $body")
    }

    @ReactMethod
    fun execWithJSONObject(action:String, obj: ReadableMap, promise:Promise){
        Log.d(NAME, "execWithJSONObject begin: $action")

        if (mAgent == null) {
            promise.reject(Exception("Not Available"))
            return
        }
        try {
            val body = JSON.toJSONBytes(obj)
            val results = mAgent.exec(action.toByteArray(),body)
            val resultsObj = JSON.parseObject(results)
            val resultsMap = Arguments.makeNativeMap(resultsObj.toMap())
            promise.resolve(resultsMap)
        }catch (e:Exception){
            Log.e(NAME, "execWithJSONObject Error: ${e.toString()}")
            promise.reject(e)
        }
        Log.d(NAME, "execWithJSONObject end: $action")
    }

    @ReactMethod
    fun execWithJSONArray(action:String, obj: ReadableMap, promise:Promise) {
        Log.d(NAME, "execWithJSONArray begin: $action")

        if (mAgent == null) {
            promise.reject(Exception("Not Available"))
            return
        }

        try {
            val body = JSON.toJSONBytes(obj)
            val results = mAgent.exec(action.toByteArray(),body)
            val resultsArr = JSON.parseArray(results)
            val resultsNativeArr = Arguments.makeNativeArray(resultsArr.toList())
            promise.resolve(resultsNativeArr)
        }catch (e:Exception){
            Log.e(NAME, "execWithJSONArray Error: ${e.toString()}")
            promise.reject(e)
        }
        Log.d(NAME, "execWithJSONArray end: $action")
    }

    override fun getName() = NAME

    companion object {
        const val NAME = "AgentNativeModule"
    }
}