{ pkgs, ... }:

let
  # The analyzer must be built with a toolchain whose go/types knows the
  # language features under test (generic methods need Go 1.27), and
  # golangci-lint parses the code with the go/types of the Go it was built
  # with. Both are therefore pinned to the same package: go_1_27 is 1.27.1 and
  # golangci-lint is 2.13.2 at the nixpkgs revision in devenv.lock. Bump by
  # running `devenv update` (and, once nixpkgs moves on, by changing the
  # attribute below).
  go = pkgs.go_1_27;

  golangci-lint = pkgs.golangci-lint.override {
    buildGo127Module = pkgs.buildGoModule.override { inherit go; };
  };
in
{
  languages.go = {
    enable = true;
    package = go;
  };

  packages = [ golangci-lint ];

  enterShell = ''
    go version
    golangci-lint version
  '';

  enterTest = ''
    go test ./...
    golangci-lint run ./...
  '';
}
