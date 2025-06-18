package com.pango.servlet

import android.Manifest
import android.content.Context
import android.content.pm.PackageManager
import android.os.Build
import androidx.core.content.ContextCompat
import gomobile.GoMobileAgent
import gomobile.Settings
import gomobile.Gomobile
import java.util.concurrent.atomic.AtomicBoolean

internal class AgentProxy() {
    @Volatile
    private var mAgent: GoMobileAgent? = null
    private val mAlready = AtomicBoolean()

    fun init(ctx: Context){
        if (!mAlready.compareAndSet(false,true)) {
            return
        }

        val settings = Settings()
        settings.configPath = ctx.dataDir.path
        settings.tempDBPath = ctx.cacheDir.path
        settings.hostName = Build.MODEL

        val appDBFile = ctx.getDatabasePath("app.db")
        settings.dbPath = appDBFile.parent

        val logger = AgentLogger()
        mAgent = Gomobile.new_(settings,logger)
    }

    fun run() {
        mAgent?.run()
    }

    fun terminate(){
        mAgent?.terminate()
    }

    fun exec(action:ByteArray,data:ByteArray):ByteArray{
        val cAgent = mAgent ?: throw Exception("Not Available")
        return cAgent.exec(action,data)
    }

    companion object {
        fun checkPermissions(context: Context):Pair<Boolean,Array<String>>{
            val permissions = arrayListOf<String>()
            var available = true

            if (ContextCompat.checkSelfPermission(context, Manifest.permission.INTERNET) != PackageManager.PERMISSION_GRANTED) {
                available = false
                permissions.add(Manifest.permission.INTERNET)
            }

            return available to permissions.toTypedArray()
        }
    }
}