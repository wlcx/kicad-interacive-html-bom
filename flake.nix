{
  description = "A static site generator for kicad projects";

  inputs.utils.url = "github:numtide/flake-utils";
  inputs.devshell = {
    url = "github:numtide/devshell";
    inputs.utils.follows = "utils";
  };

  outputs = {
    self,
    nixpkgs,
    utils,
    devshell,
  }:
    utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {
        inherit system;
        overlays = [devshell.overlay];
      };
    in rec {
      packages.default = pkgs.buildGoModule {
        name = "kicad-site-generator";
        src = self;
        vendorSha256 = "d1mQwSPnZDAMvBAUEdi+Sogr49LLzZ3TXZoztZU6sjo=";

        # Inject the git version
        ldflags = ''
          -X main.version=${if self ? rev then self.rev else "dirty"}
        '';
      };

      apps.default = utils.lib.mkApp {drv = packages.default;};

      devShells.default =
        pkgs.devshell.mkShell {packages = with pkgs; [go gopls];};
      formatter = pkgs.alejandra;
    });
}
