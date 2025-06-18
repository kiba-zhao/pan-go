package com.pango.servlet

import android.view.View
import com.facebook.react.ReactPackage
import com.facebook.react.bridge.NativeModule
import com.facebook.react.bridge.ReactApplicationContext
import com.facebook.react.uimanager.ReactShadowNode
import com.facebook.react.uimanager.ViewManager
import gomobile.GoMobileAgent

internal class ReactNativePackage(agentProxy: AgentProxy): ReactPackage {
    private val mAgentProxy = agentProxy

    override fun createViewManagers(
        reactContext: ReactApplicationContext
    ): MutableList<ViewManager<View, ReactShadowNode<*>>> = mutableListOf()

    override fun createNativeModules(
        reactContext: ReactApplicationContext
    ): MutableList<NativeModule> = listOf(AgentNativeModule(reactContext,mAgentProxy)).toMutableList()
}