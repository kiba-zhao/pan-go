plugins {
    id("com.android.library")
    id("org.jetbrains.kotlin.android")
}

android {

    ndkVersion = rootProject.ext.get("ndkVersion") as String
    buildToolsVersion = rootProject.ext.get("buildToolsVersion") as String
    compileSdk = rootProject.ext.get("compileSdkVersion") as Int

    namespace = "com.pango.servlet"

    defaultConfig {
        minSdk = rootProject.ext.get("minSdkVersion") as Int
        lint.targetSdk = rootProject.ext.get("targetSdkVersion") as Int
        //testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
        consumerProguardFiles("consumer-rules.pro")
    }

    buildTypes {
        release {
            isMinifyEnabled = false
            proguardFiles(
                getDefaultProguardFile("proguard-android-optimize.txt"),
                "proguard-rules.pro"
            )
        }
    }
   /*
   compileOptions {
        sourceCompatibility = JavaVersion.VERSION_11
        targetCompatibility = JavaVersion.VERSION_11
    }
    kotlinOptions {
        jvmTarget = "11"
    }
    */

}

dependencies {
    implementation(fileTree(mapOf("dir" to "libs")))

    implementation("com.facebook.react:react-android")

    implementation("androidx.core:core-ktx:1.13.1")
    implementation("androidx.appcompat:appcompat:1.7.0")
    /*
implementation("com.google.android.material:material:1.10.0")
testImplementation("junit:junit:4.13.2")
androidTestImplementation("androidx.test.ext:junit:1.1.5")
androidTestImplementation("androidx.test.espresso:espresso-core:3.5.1")
*/
}