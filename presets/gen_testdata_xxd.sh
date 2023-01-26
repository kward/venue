#! /bin/bash

srcDir="${HOME}/Documents/D-Show/User Data/Effect Presets/testdata/"

if [ ! -d "${srcDir}" ]; then
  echo "missing destination directory '${destDir}'"
  exit 1
fi
cd "${srcDir}"

gen() {
  d=$1
  f=$2

  xxd "${f}" >"${f}.xxd"
  echo "$(date) + ${d}/${f}.xxd"
}

while true; do
  ls |while read d; do
    pushd "${d}" >/dev/null

    # Generate new .xxd files.
    ls *.ich 2>/dev/null |while read f; do
      [ -r "${f}.xxd" ] || gen "${d}" "${f}"
      [ "${f}.xxd" -nt "${f}" ] || gen "${d}" "${f}"
    done

    # Clean old .xxd files.
    ls *.xxd 2>/dev/null |while read f; do
      if [ ! -r "${f/.xxd/}" ]; then
       rm "${f}"
       echo "$(date) - ${d}/${f}"
      fi
    done

    popd >/dev/null
  done

  sleep 3
done
