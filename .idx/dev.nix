{ pkgs, ... }: {
  channel = "stable-24.05"; # or "unstable"
  packages = [
    pkgs.go
  ];
  env = {};
  idx = {
    previews = {
      enable = false;
    };workspace = {
      onCreate = {
        go-install = "go get";        
        default.openFiles = [ ".idx/dev.nix" "README.md" ];
      };
    };
  };
}
