{
  description = "DDoS botnet written in Go";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-22.11";
  };
  outputs = { self, nixpkgs }:
    let
      pkgs = nixpkgs.legacyPackages.x86_64-linux;
      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";
      version = builtins.substring 0 8 lastModifiedDate;
    in
    {
      formatter.x86_64-linux = pkgs.nixpkgs-fmt;
      packages.x86_64-linux.apisrv = pkgs.buildGoModule {
        pname = "apisrv";
        src = ./.;
        inherit version;
        vendorSha256 = "sha256-03NpYZL4T9Iz6Ts9+bMl3dMAqNQ09aZt8xjkSqNU+EI=";
      };
      packages.x86_64-linux.daemon = pkgs.buildGoModule {
        pname = "daemon";
        src = ./.;
        inherit version;
        vendorSha256 = "sha256-03NpYZL4T9Iz6Ts9+bMl3dMAqNQ09aZt8xjkSqNU+EI=";
      };
      packages.x86_64-linux.botctl = pkgs.buildGoModule {
        pname = "botctl";
        src = ./.;
        inherit version;
        vendorSha256 = "sha256-03NpYZL4T9Iz6Ts9+bMl3dMAqNQ09aZt8xjkSqNU+EI=";
      };
      devShells.x86_64-linux.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          gopls
          protobuf
          protoc-gen-go
          protoc-gen-go-grpc
        ];
      };
    };
}
