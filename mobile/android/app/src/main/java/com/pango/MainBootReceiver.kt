package com.pango

import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.util.Log
import androidx.core.content.ContextCompat

class MainBootReceiver : BroadcastReceiver(){
    override fun onReceive(context: Context, intent: Intent) {
        Log.d("MainBootReceiver", "onReceive begin")
        if (Intent.ACTION_BOOT_COMPLETED == intent.action) {
            startMainService(context)
        }
        Log.d("MainBootReceiver", "onReceive end")
    }

    private fun startMainService(context: Context){
        Log.d("MainBootReceiver", "startMainService begin")

        val (available,_) = MainService.checkPermissions(context)
        if (!available){
            Log.w("MainBootReceiver", "startMainService Not Available")
            return
        }

        val serviceIntent = Intent(context, MainService::class.java)
        ContextCompat.startForegroundService(context, serviceIntent)
        Log.d("MainBootReceiver", "startMainService end")
    }
}