if [ $# -eq 0 ]; then
    echo "You must pass the version number."
    exit 1
fi

go build cmd/main.go
git commit -m "release: version $1"
git tag -a "$1" -m "Version $1"
