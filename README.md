# PostgreScrutiniser

PostgreSQL configuration management tool, inspired by [MySQLTuner](https://github.com/major/MySQLTuner-perl).

This is the code repository for my final software development bachelor project.

## Installation

The frontend and backend sides can be installed separately. The application has been created and fully tested on AlmaLinux 8.7. This project assumes that passwords inside `/etc/shadow` are hashed with the "SHA512" algorithm. 

### Backend setup

To setup the backend part of this project, change directory to `backend/` and run `setup.sh`:
```
cd backend/
./setup.sh
```

This will create all the necessary execution related files as well as a systemd service file, allowing you to control it using `systemctl`. Make sure that correct PostgreSQL user and password names are provided as that's what what will be used by our application to communicate with PostgreSQL.

The setup script will also output login details for our main user. If you end up losing it, simply change the password with `passwd postgrescrutiniser` and use the new password.

After setup is complete, feel free to remove this repository.

### Frontend setup

To setup the frontend part of this project, build it for production with:
```
cd frontend/
npm install
npm run build
```

Depending on where the frontend part of this project is being setup, the way it should be served will differ. Easiest option is to install `http-server` package and then running it inside `frontend/dist/` directory:
```
npm install -g http-server
http-server
```

After that enter the `http` address that's given below `Available on:` part of the output.

Note that only the files inside `dist/` directory would then be of any importance. The rest can be removed.

## Usage guide

