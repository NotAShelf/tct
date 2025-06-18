{
  lib,
  buildGoModule,
}: let
  fs = lib.fileset;
  s = ../.;
in
  buildGoModule {
    pname = "tct";
    version = "0.1.0";

    src = fs.toSource {
      root = s;
      fileset = fs.unions [
        (fs.fileFilter (file: builtins.any file.hasExt ["go"]) s)
        ../go.mod
        ../go.sum
      ];
    };

    vendorHash = "sha256-m5mBubfbXXqXKsygF5j7cHEY+bXhAMcXUts5KBKoLzM=";

    ldflags = ["-s" "-w"];
  }
