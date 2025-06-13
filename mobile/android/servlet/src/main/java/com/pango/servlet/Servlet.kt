package com.pango.servlet


import android.content.ContextWrapper
import android.os.Build
import gomobile.GoMobileAgent
import gomobile.Settings
import gomobile.Gomobile

class Servlet(ctx: ContextWrapper) {

    private val mAgent: GoMobileAgent? = newGoMobileAgent(ctx)
    val reactPackage = ReactNativePackage(mAgent)

    fun start(){
        mAgent?.run()
    }

    fun stop(){
        mAgent?.terminate()
    }

    companion object {
        private fun newGoMobileAgent(ctx: ContextWrapper):GoMobileAgent?{

            val mountDir = ctx.getExternalFilesDir(null) ?: return null
            val settings = Settings()
            settings.mountPath = mountDir.path+"/"+ctx.packageName

            val appDBFile = ctx.getDatabasePath("app.db")
            settings.dbPath = appDBFile.parent

            settings.configPath = ctx.dataDir.path
            settings.tempDBPath = ctx.cacheDir.path
            settings.hostName = Build.MODEL



            val logger = AgentLogger()
            return Gomobile.new_(settings,logger)
        }
    }

}