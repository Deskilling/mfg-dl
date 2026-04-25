{
  description = "Mfg-dl";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "mfg-dl";
          version = "0.0.1";
          src = ./.;

          vendorHash = "sha256-BrXg35DXKhqdVVyEyiL0uyUm41zwDajNtdcEavHaVzg=";
          subPackages = [ "cmd/mfg-dl" ];

          ldflags = [
            "-s"
            "-w"
          ];

          trimpath = true;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            gnumake
            go
          ];
        };
      }
    );
}
