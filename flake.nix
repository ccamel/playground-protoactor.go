{
  description = "playground-protoactor.go development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";
  };

  outputs =
    { nixpkgs, ... }:
    let
      supportedSystems = [
        "aarch64-darwin"
        "x86_64-linux"
        "aarch64-linux"
      ];

      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
    in
    {
      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.mkShell {
            packages = [
              pkgs.bash-language-server
              pkgs.git
              pkgs.go_1_26
              pkgs.gofumpt
              pkgs.golangci-lint
              pkgs.gopls
              pkgs.markdownlint-cli2
              pkgs.yaml-language-server
            ];

            shellHook = ''
              echo "playground-protoactor.go development environment loaded"
              echo "Go: $(go version)"
            '';
          };
        }
      );
    };
}
