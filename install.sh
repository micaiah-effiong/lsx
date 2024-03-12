## build
mkdir -p build

go build -o build main.go

chmod +x build/main

## copy to user path
LSX_PATH="$HOME/.lsx"
LSX_FN_SCRIPT="$LSX_PATH/lsx.sh"

mkdir -p $LSX_PATH

mv "build/main" $LSX_PATH
cp "lsx.sh" $LSX_PATH

echo "#################"
echo "$(ls $LSX_PATH)"
echo "#################"

## add lsx.sh to user ...rc file
echo $(which $SHELL)
echo "source \$HOME/.lsx/lsx.sh" >> "$HOME/.zshrc"
# echo "source $LSX_FN_SCRIPT" >> "\$HOME/.bashrc"

