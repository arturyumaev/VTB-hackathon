import React from 'react'
import { Button } from 'antd'
import Fingerprint2 from '@fingerprintjs/fingerprintjs'

export default function SuspectActivity() {
  const handleClick = () => {
    var options = { excludes: {
      touchSupport: true,
      webdriver: true,
      fonts: true,
      audio: true,
      hasLiedBrowser: true,
      hasLiedOs: true,
      hasLiedResolution: true,
      hasLiedLanguages: true,
      webgl: true,
      canvas: true,
      cpuClass: true,
      plugins: true,
      enumerateDevices: true,
      fontsFlash: true,
      adBlock: true,
      doNotTrack: true,
      addBehavior: true,
      openDatabase: true,
      localStorage: true,
      sessionStorage: true,
      pixelRatio: true,
      colorDepth: true,
      indexedDb: true,
    }}

    if (window.requestIdleCallback) {
      requestIdleCallback(function () {
          Fingerprint2.get(options, function (components) {
            var values = components.map(c => c.value)
            var fingerprintHash = Fingerprint2.x64hash128(values.join(''), 31)

            console.log(fingerprintHash)
          })
      })
    } else {
        setTimeout(function () {
          Fingerprint2.get(options, function (components) {
            var values = components.map(c => c.value)
            var fingerprintHash = Fingerprint2.x64hash128(values.join(''), 31)

            console.log(fingerprintHash)
          })  
        }, 500)
    }
  }

  return (
    <div>
      <Button onClick={handleClick}>Click</Button>
    </div>
  )
}
