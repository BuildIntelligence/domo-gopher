# Getting Started
- Create a .env file (example contents below) with values for Domo client auth in the root directory for the tests to work. (i.e. same directory as this README.md file)
```
# Domo Environment
DOMO_CLIENT_ID=abcHeresTheClientId
DOMO_SECRET=shhItsASecretValue!
DATA_SCOPE=true
USER_SCOPE=true
AUDIT_SCOPE=true
DASHBOARD_SCOPE=true
```

# TODO:
- [] improve README
- [x] improve auth scope configuration to include scope in the url auth params based on input flags
- [x] Dataset API wrapper methods
- [x] Stream API wrapper methods
- [] User API wrapper methods
- [] Group API wrapper methods
- [] Page API wrapper methods
- [] Go Modules for dependency management

