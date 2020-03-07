#!/bin/bash

# 1. Install chromedriver
# 2. Install geckodriver

# ChromeDriver options
chromeDriverReleaseURL=https://chromedriver.storage.googleapis.com/LATEST_RELEASE
chromeDriverDefaultVersion=80.0.3987.105
chromeDriverName=chromedriver

# GeckoDriver options
geckoDriverDefaultVersion=v0.26.0
geckoDriverName=geckodriver

function retrieveChromeDriverVersion() {
    version=$(curl -s "${chromeDriverReleaseURL}")
    if [ -z "$version" ]; then
        version=${chromeDriverDefaultVersion}
    fi
}

function retrieveGeckoDriverVersion() {
    version=${geckoDriverDefaultVersion}
}

retrieveChromeDriverVersion
echo "Install ChromeDriver. Version=${version}"
declare -a platforms=("linux64" "mac64")
for platform in "${platforms[@]}"
do
     name="${chromeDriverName}_${platform}"
     curl -s https://chromedriver.storage.googleapis.com/$version/$name.zip -O
     unzip -q -o ${name}.zip
     rm ${name}.zip
     if [[ -f "chromedriver" ]]; then
       mv chromedriver ${name}
       chmod a+x ${name}
     fi
done

retrieveGeckoDriverVersion
echo "Install GeckoDriver. Version=${version}"
declare -a platforms=("linux64" "macos")
for platform in "${platforms[@]}"
do
   name="${geckoDriverName}_${platform}"
   curl -Ls https://github.com/mozilla/geckodriver/releases/download/${version}/geckodriver-${version}-${platform}.tar.gz | tar xz
   if [ -f "geckodriver" ]; then
      mv geckodriver ${name}
      chmod a+x ${name}
   fi
done


