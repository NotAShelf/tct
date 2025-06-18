{
  mkShell,
  go,
  gopls,
  delve,
}:
mkShell {
  name = "go";
  packages = [
    delve
    go
    gopls
  ];
}
