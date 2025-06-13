package com.pango

import android.content.ComponentName
import android.content.Context
import android.content.Intent
import android.content.ServiceConnection
import android.os.Bundle
import android.os.IBinder
import android.util.Log
import androidx.core.app.ActivityCompat
import androidx.core.content.ContextCompat
import com.facebook.react.ReactActivity
import com.facebook.react.ReactActivityDelegate
import com.facebook.react.defaults.DefaultNewArchitectureEntryPoint.fabricEnabled
import com.facebook.react.defaults.DefaultReactActivityDelegate

class MainActivity : ReactActivity() {

    private var mainServiceAlready = false

  /**
   * Returns the name of the main component registered from JavaScript. This is used to schedule
   * rendering of the component.
   */
  override fun getMainComponentName(): String = "pango"

  /**
   * Returns the instance of the [ReactActivityDelegate]. We use [DefaultReactActivityDelegate]
   * which allows you to enable New Architecture with a single boolean flags [fabricEnabled]
   */
  override fun createReactActivityDelegate(): ReactActivityDelegate =
      DefaultReactActivityDelegate(this, mainComponentName, fabricEnabled)

    override fun onCreate(savedInstanceState: Bundle?) {
        Log.d("MainActivity", "onCreate begin")
        super.onCreate(savedInstanceState)

        startMainService(true)
        Log.d("MainActivity", "onCreate end")
    }

    override fun onDestroy() {
        Log.d("MainActivity", "onDestroy begin")
        super.onDestroy()

        stopMainService()
        Log.d("MainActivity", "onDestroy end")
    }

    override fun onRequestPermissionsResult(
        requestCode: Int,
        permissions: Array<out String>,
        grantResults: IntArray
    ) {
        Log.d("MainActivity", "onRequestPermissionsResult begin $requestCode")
        super.onRequestPermissionsResult(requestCode, permissions, grantResults)
        if (requestCode == REQUEST_MAIN_SERVICE_PERMISSION){
            startMainService()
        }
        Log.d("MainActivity", "onRequestPermissionsResult end $requestCode")
    }

    private fun serviceIntent():Intent{
        return Intent(this, MainService::class.java)
    }

    private fun startMainService(tryRequestPermissions:Boolean=false){

        Log.d("MainActivity", "startMainService begin $tryRequestPermissions")

        // check permissions
        val (available,permissions) = MainService.checkPermissions(this)
        Log.d("MainActivity", "startMainService checkPermissions $available, ${permissions.joinToString()}")
        if (permissions.isNotEmpty()){
            if (tryRequestPermissions){
//                if (ActivityCompat.shouldShowRequestPermissionRationale(this,Manifest.permission.POST_NOTIFICATIONS)){
//
//                }
                ActivityCompat.requestPermissions(this, permissions,REQUEST_MAIN_SERVICE_PERMISSION)
                Log.d("MainActivity", "startMainService Try Request Permissions $tryRequestPermissions")
                return
            }
            if(!available){
                Log.w("MainActivity", "startMainService Not Available $tryRequestPermissions")
                return
            }
        }
        //
        ContextCompat.startForegroundService(this, serviceIntent())
        mainServiceAlready = true
        Log.d("MainActivity", "startMainService end $tryRequestPermissions")
    }

    private fun stopMainService(){
        Log.d("MainActivity", "stopMainService begin")
        if (mainServiceAlready) {
            stopService(serviceIntent())
        }
        Log.d("MainActivity", "stopMainService end")
    }

    companion object {
        const val REQUEST_MAIN_SERVICE_PERMISSION = 1
    }
}
