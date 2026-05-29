foreach ($line in Get-Content "$PSScriptRoot/../.env") {
	if ($line -match '^\s*#' -or $line -match '^\s*$') {
		continue
	}

	$key, $value = $line -split '=', 2
	[System.Environment]::SetEnvironmentVariable($key, $value)
}