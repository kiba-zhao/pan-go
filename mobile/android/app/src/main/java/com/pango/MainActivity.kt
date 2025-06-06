package com.pango

import android.content.ComponentName
import android.content.Context
import android.content.Intent
import android.content.ServiceConnection
import android.os.Bundle
import android.os.IBinder
import android.util.Log
import androidx.core.app.ActivityCompat
import com.facebook.react.ReactActivity
import com.facebook.react.ReactActivityDelegate
import com.facebook.react.defaults.DefaultNewArchitectureEntryPoint.fabricEnabled
import com.facebook.react.defaults.DefaultReactActivityDelegate

class MainActivity : ReactActivity() {

    private lateinit var mainServiceBinder: MainService.MainServiceBinder
    private lateinit var mainServiceConnection: ServiceConnection
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

        bindMainService(true)
        Log.d("MainActivity", "onCreate end")
    }

    override fun onDestroy() {
        Log.d("MainActivity", "onDestroy begin")
        super.onDestroy()

        unbindMainService()
        Log.d("MainActivity", "onDestroy end")
    }

    fun getMainServiceBinder(): MainService.MainServiceBinder = mainServiceBinder

    override fun onRequestPermissionsResult(
        requestCode: Int,
        permissions: Array<out String>,
        grantResults: IntArray
    ) {
        Log.d("MainActivity", "onRequestPermissionsResult begin $requestCode")
        super.onRequestPermissionsResult(requestCode, permissions, grantResults)
        if (requestCode == REQUEST_MAIN_SERVICE_PERMISSION){
            bindMainService()
        }
        Log.d("MainActivity", "onRequestPermissionsResult end $requestCode")
    }

    private fun bindMainService(tryRequestPermissions:Boolean=false){

        Log.d("MainActivity", "bindMainService begin $tryRequestPermissions")

        // check permissions
        val (available,permissions) = MainService.checkPermissions(this)
        Log.d("MainActivity", "bindMainService checkPermissions $available, ${permissions.joinToString()}")
        if (permissions.isNotEmpty()){
            if (tryRequestPermissions){
//                if (ActivityCompat.shouldShowRequestPermissionRationale(this,Manifest.permission.POST_NOTIFICATIONS)){
//
//                }
                ActivityCompat.requestPermissions(this, permissions,REQUEST_MAIN_SERVICE_PERMISSION)
                Log.d("MainActivity", "bindMainService Try Request Permissions $tryRequestPermissions")
                return
            }
            if(!available){
                Log.w("MainActivity", "bindMainService Not Available $tryRequestPermissions")
                return
            }
        }
        //

        mainServiceConnection =
            object : ServiceConnection {
                override fun onServiceConnected(name: ComponentName, service: IBinder) {
                    Log.d("MainActivity", "ServiceConnection.onServiceConnected before")
                    mainServiceBinder = service as MainService.MainServiceBinder
                    Log.d("MainActivity", "ServiceConnection.onServiceConnected end")
                }

                override fun onServiceDisconnected(name: ComponentName) {}
            }
        val serviceIntent = Intent(this, MainService::class.java)

        bindService(serviceIntent, mainServiceConnection, Context.BIND_AUTO_CREATE)
        mainServiceAlready = true
        Log.d("MainActivity", "bindMainService end $tryRequestPermissions")
    }

    private fun unbindMainService(){
        Log.d("MainActivity", "unbindMainService begin")
        if (!mainServiceAlready){
            Log.d("MainActivity", "unbindMainService Not Already")
            return
        }
        unbindService(mainServiceConnection)
        Log.d("MainActivity", "unbindMainService end")
    }

    companion object {
        const val REQUEST_MAIN_SERVICE_PERMISSION = 1
    }
}
