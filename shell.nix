{pkgs ? import <nixpkgs> {}}:
pkgs.mkShell {
  buildInputs = with pkgs; [
    go
  ];

  shellHook = ''
    export GOPATH=$PWD/.gopath
    export GOBIN=$PWD/.gopath/bin
    export PATH=$GOBIN:$PATH

    go install github.com/bootdotdev/bootdev@latest

    export SHELL=$(which zsh)
    exec $SHELL
  '';
}
