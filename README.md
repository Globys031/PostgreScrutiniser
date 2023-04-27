# PostgreScrutiniser

PostgreSQL configuration management tool, inspired by [MySQLTuner](https://github.com/major/MySQLTuner-perl).

This is the code repository for my final software development bachelor project.

## Installation

The frontend and backend sides can be installed separately.

### Backend setup

To setup the backend part of this project, change directory to `backend/` and run `setup.sh`:
```
cd backend/
./setup.sh
```

This will create all the necessary execution related files as well as a systemd service file, allowing you to control it using `systemctl`.

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