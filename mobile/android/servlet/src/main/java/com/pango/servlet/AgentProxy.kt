package com.pango.servlet

import android.Manifest
import android.content.Context
import android.content.pm.PackageManager
import android.net.ConnectivityManager
import android.net.NetworkCapabilities
import android.net.wifi.WifiManager
import android.os.Build
import android.util.Log
import androidx.core.content.ContextCompat
import gomobile.GoMobileAgent
import gomobile.Gomobile
import gomobile.Settings
import java.net.InetAddress
import java.net.NetworkInterface
import java.util.Collections
import java.util.concurrent.atomic.AtomicBoolean


internal class AgentProxy() {
    @Volatile
    private var mAgent: GoMobileAgent? = null
    private val mAlready = AtomicBoolean()

    fun sampleInterfaces(){
        Log.d("SampleTest","NetworkInterface begin")
        val interfaces: List<NetworkInterface> =
            Collections.list(NetworkInterface.getNetworkInterfaces())
        for ( iface in interfaces) {
            Log.d("SampleTest",
                "NetworkInterface range: name="+iface.name
                        +", displayName="+iface.displayName+
                        ", mtu="+iface.mtu+
                        ", isMulitcast="+iface.supportsMulticast()
                        +", isUp="+iface.isUp
                        +", isLoopback="+iface.isLoopback
                        +", isVirtual="+iface.isVirtual
                        +", isPonitToPonit="+iface.isPointToPoint
            )


            for (netAddr in iface.inetAddresses){
                Log.d("SampleTest",
                    "NetworkInterface inetAddresses: name="+iface.name
                            +", hostAddr="+netAddr.hostAddress
                )
            }

        }
    }

    fun sampleWfiAddr(ctx: Context){
        val cm = ctx.getSystemService(Context.CONNECTIVITY_SERVICE) as ConnectivityManager?
      
        val network = cm!!.activeNetwork
        val nc = cm.getNetworkCapabilities(network)
        if (nc != null && nc.hasTransport(NetworkCapabilities.TRANSPORT_WIFI)) {
            Log.d("SampleTest","with ConnectivityManager begin")
            val linkProperties = cm.getLinkProperties(network)
            if (linkProperties != null){
                val iface = NetworkInterface.getByName(linkProperties.interfaceName)
                for (addr in iface.interfaceAddresses){
                    Log.d("SampleTest","with ConnectivityManager: addr="+addr.address.hostAddress+"/"+addr.networkPrefixLength
                            +" ,name="+iface.name
                            +", isLinkLocalAddress="+addr.address.isLinkLocalAddress
                    )
                }

            }
        }


    }

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

        // sample test android api
        sampleWfiAddr(ctx)



    }

    fun run() {
        //mAgent?.run()
    }

    fun terminate(){
        //mAgent?.terminate()
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