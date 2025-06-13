package com.pango.servlet

import android.util.Log
import gomobile.Logger

class AgentLogger :Logger{
    override fun debug(tag:String,msg:String){
        Log.d("$LOG_TAG","$tag:$msg")
    }

    override fun info(tag:String,msg:String) {
        Log.i("$LOG_TAG","$tag:$msg")
    }

    override fun error(tag:String,msg:String) {
        Log.e("$LOG_TAG","$tag:$msg")
    }

    override fun warn(tag:String,msg:String) {
        Log.w("$LOG_TAG","$tag:$msg")
    }
    companion object {
        const val LOG_TAG = "GoLogger"
    }
}