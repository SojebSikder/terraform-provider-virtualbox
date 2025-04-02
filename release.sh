#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status

VERSION="0.0.1"  # Update this to match your version
tf_provider_name="terraform-provider-virtualbox"
build_dir="build"

# Create the build directory if it doesn't exist
mkdir -p $build_dir

# List of target platforms
platforms=(
  "linux amd64"
  "windows amd64"
  "darwin amd64"
)

# Build and package for each platform
for platform in "${platforms[@]}"; do
  set -- $platform
  GOOS=$1
  GOARCH=$2
  output_name="$tf_provider_name_v$VERSION"
  
  if [ "$GOOS" == "windows" ]; then
    output_name+=".exe"
  fi
  
  echo "Building for $GOOS/$GOARCH..."
  env GOOS=$GOOS GOARCH=$GOARCH go build -o $build_dir/$output_name

  zip_name="$tf_provider_name_${VERSION}_${GOOS}_${GOARCH}.zip"
  echo "Zipping $zip_name..."
  zip -j "$build_dir/$zip_name" "$build_dir/$output_name"
done

# Generate SHA256SUMS file
echo "Generating SHA256SUMS..."
shasum -a 256 $build_dir/*.zip > "$build_dir/${tf_provider_name}_${VERSION}_SHA256SUMS"

# Sign the SHA256SUMS file
echo "Signing SHA256SUMS..."
gpg --detach-sign "$build_dir/${tf_provider_name}_${VERSION}_SHA256SUMS"

echo "Release files prepared successfully in $build_dir"