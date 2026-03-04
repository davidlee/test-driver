{
  description = "A python flake-parts shell for im";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    devshell.url = "github:numtide/devshell";
  };

  outputs = inputs @ {flake-parts, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} {
      imports = [
        inputs.devshell.flakeModule
      ];

      systems = [
        "x86_64-linux"
        "aarch64-darwin"
      ];

      perSystem = {pkgs, ...}: {
        devshells.default = {
          packages = with pkgs;
            [
              #####################
              # language support
              #####################

              ## Go
              go
              gomarkdoc
              golangci-lint
              # go # compiler
              # gopls # language server
              # gomacro
              # gofumpt # strict formatter
              # golangci-lint # linter

              # python
              uv
              python3Packages.python-lsp-server
              python3Packages.python-lsp-ruff
              tree-sitter
              tree-sitter-grammars.tree-sitter-python
              pyright

              # js / ts
              # nodejs_latest
              bun

              ## diagrams
              d2
              graphviz

              ## utils
              just
              # watchexec
            ]
            ++ lib.optionals stdenv.isLinux [];

          commands = [
            {
              name = "sdr";
              help = "spec-driver $@";
              command = ''
                spec-driver $@
              '';
            }
          ];
        };
      };
    };
}
