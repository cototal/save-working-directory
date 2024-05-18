# Windows

Use in conjunction with PowerShell:

```powershell
# Depends on the save-working-directory being installed to path
function Set-GoWorkingDirectory($Project = "default") {
    save-working-directory -s $Project
}
Set-Alias swd Set-GoWorkingDirectory

# Depends on the save-working-directory being installed to path
function Get-GoWorkingDirectory($Project = "default") {
    $output = save-working-directory -l $Project
    if ($LASTEXITCODE -ne 1) {
        Set-Location -Path $output
    }
}
Set-Alias lwd Get-GoWorkingDirectory

# Depends on the save-working-directory being installed to path
function Remove-GoWorkingDirectory($Project) {
    save-working-directory -d $Project
}
Set-Alias dwd Remove-GoWorkingDirectory

# Depends on the save-working-directory being installed to path
function Show-GoWorkingDirectories() {
    save-working-directory --list
}
Set-Alias wdlist Show-GoWorkingDirectories
```

# Linux

Save to path as `swd`. Then add alias with Bash:

```bash
function lwd() {
    local name="${1:-default}"
    local path=$(swd -l $name)
    if [ $? -eq 0 ]; then
        cd $path
    fi
}
```