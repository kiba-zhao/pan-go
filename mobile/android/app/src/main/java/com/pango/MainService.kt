package com.pango

import android.Manifest
import android.app.Service
import android.content.Context
import android.content.Intent
import android.content.pm.PackageManager
import android.content.pm.ServiceInfo
import android.os.Binder
import android.os.Build
import android.os.IBinder
import android.util.Log
import androidx.core.app.NotificationCompat
import androidx.core.app.ServiceCompat
import androidx.core.content.ContextCompat

class MainService : Service() {
    inner class MainServiceBinder : Binder() {
        fun getService(): MainService = this@MainService
    }

    override fun onBind(intent: Intent): IBinder {

        Log.d("MainService", "onBind begin")
        val binder = MainServiceBinder()
        Log.d("MainService", "onBind end")
        return binder

        //        return MainServiceBinder()
    }

    override fun onCreate() {
        // Log
        Log.d("MainService", "onCreate begin")
        //

        super.onCreate()

        startForegroundMainService()
        startGoBinder()

        Log.d("MainService", "onCreate end")
    }

    override fun onDestroy() {
        Log.d("MainService", "onDestroy begin")
        super.onDestroy()

        stopGoBinder()
        Log.d("MainService","onDestroy end")
    }

    private fun startForegroundMainService(){
        Log.d("MainService", "startForegroundMainService begin")
        val notification = NotificationCompat.Builder(this, MAIN_SERVICE_CHANNEL_ID)
            .setSmallIcon(R.mipmap.ic_launcher)
            .setContentTitle("TestApp")
            .setContentText("MainService is running")
            .build()

        ServiceCompat.startForeground(
            this,
            MAIN_SERVICE_ID,
            notification,
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.Q) {
                ServiceInfo.FOREGROUND_SERVICE_TYPE_DATA_SYNC
            } else {
                0
            }
        )
        Log.d("MainService", "startForegroundMainService end")
    }

    private fun startGoBinder(){
        Log.d("MainService", "startGoBinder begin")
        // TODO: impelement
        Log.d("MainService", "startGoBinder end")
    }

    private fun stopGoBinder(){
        Log.d("MainService", "stopGoBinder begin")
        // TODO: impelement
        Log.d("MainService", "stopGoBinder end")
    }

    companion object {
        const val MAIN_SERVICE_ID = 1
        const val MAIN_SERVICE_CHANNEL_ID = "MainServiceForegroundChannel"
        fun checkPermissions(context: Context):Pair<Boolean,Array<String>>{
            val permissions = arrayListOf<String>()
            var available = true
            Log.d("MainService","checkPermissions POST_NOTIFICATIONS ${Build.VERSION.SDK_INT}")
            if (
                Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU
                && (ContextCompat.checkSelfPermission(context,
                    Manifest.permission.POST_NOTIFICATIONS) != PackageManager.PERMISSION_GRANTED)
            ) {
                permissions.add(Manifest.permission.POST_NOTIFICATIONS)
            }

            if (
                Build.VERSION.SDK_INT >= Build.VERSION_CODES.Q
                && (ContextCompat.checkSelfPermission(context,
                    Manifest.permission.FOREGROUND_SERVICE) != PackageManager.PERMISSION_GRANTED)
            ) {
                available = false
                permissions.add(Manifest.permission.FOREGROUND_SERVICE)
            }

            if (
                Build.VERSION.SDK_INT >= Build.VERSION_CODES.UPSIDE_DOWN_CAKE
                && (ContextCompat.checkSelfPermission(context,
                    Manifest.permission.FOREGROUND_SERVICE_DATA_SYNC) != PackageManager.PERMISSION_GRANTED)
            ) {
                available = false
                permissions.add(Manifest.permission.FOREGROUND_SERVICE_DATA_SYNC)
            }

            return available to permissions.toTypedArray()
        }
    }
}