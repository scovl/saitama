#!/bin/bash

# Array of operating systems and architectures
platforms=("windows/amd64" "darwin/amd64" "linux/amd64")

# Binary name
binary_name="saitama"

# Loop through each platform
for platform in "${platforms[@]}"
do
    # Split the platform string into OS and architecture
    platform_split=(${platform//\// })
    os=${platform_split[0]}
    arch=${platform_split[1]}

    # Set output file name based on OS
    output_name=$binary_name'-'$os'-'$arch
    if [ $os = "windows" ]; then
        output_name+='.exe'
    fi

    # Build the binary
    env GOOS=$os GOARCH=$arch go build -o build/$output_name main.go
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done

echo 'Build completed successfully!'
