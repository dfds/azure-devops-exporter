
$env:AWS_PROFILE="saml"
Write-Output "Login with your email and password and select Account: dfds (454234050858) / ADFS-Admin"
saml2aws login --force
$Session = aws sts assume-role --role-arn "arn:aws:iam::738063116313:role/OrgRole" --role-session-name AWSCLI-Session --region eu-central-1 | ConvertFrom-Json
Write-Output "export AWS_ACCESS_KEY_ID=$($Session.Credentials.AccessKeyId)"
Write-Output "export AWS_SECRET_ACCESS_KEY=$($Session.Credentials.SecretAccessKey)"
Write-Output "export AWS_SESSION_TOKEN=$($Session.Credentials.SessionToken)"