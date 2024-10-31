# Minimal Viable Product (Sensor Data Acquistion)

This project adopts a modular design and a different architectural approach, serving as a proof of concept applicable to both current and future software products at Mesa. It will ingest sensor data from external sources, store and process that data to generate reports, and display the information in real-time.
    -Golang programming language
    -Apache Kafka message broker
    -QuestDB database
    -Svelte.js and Sveltekit
    -Tailwind CSS 
    -Websocket protocol


## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)


## Features

- Live display of sensor data on the dashboard over websocket with tabular view and chart view
- Utilize QuestDB high performance SQL to generate large data set report including historical, Minimum values, Maximum values, and Average values
- Utilizing docker to setup the development environment 


## Requirements

1. Windows Subsystem Linux (WSL), this is needed to create and build Docker images and containers if you are running Windows OS only.
2. Docker Desktop at: https://www.docker.com/products/docker-desktop/
3. Node.js with NPM at: https://nodejs.org/en
4. A development environment if you would like to continue developing the project, VS Code IDE at: https://code.visualstudio.com/
5. Git for windows if your OS is Windows: https://git-scm.com/downloads/win
6. Git for MAC if your OS is macOS: https://git-scm.com/downloads/mac
7. Git for Linux and Unix if your OS is such: https://git-scm.com/downloads/linux


## Installation

1. **Install Windows Subsystem Linux (WSL):**
    a. If you are running Windows, please continue, if you are using MAC or Linux OS, please skip this step and move to Step 2: Install Docker Desktop
    b.	Go to Windows Store, search for “ubuntu” and choose WSL Ubuntu 24.04 LTS, this is the latest long term support version at the time of this writing.
        ![Ubuntu WSL Installation](image-1.png)

    c.	Install the Ubuntu WSL from the Windows Store

2. **Install Docker Desktop:**
    a. Please go to: https://www.docker.com/products/docker-desktop/ and download Docker Desktop for your operating system
    b. Ensure you set your machine to enable Virtualization.  To verify on Windows, open Task Manager, under Performance tab, you should see "Virtualization: Enabled" as shown below:
       ![System Virtualization](image-2.png)
    c. Install Docker Desktop using the file you just downloaded.  On Windows OS, when ask if you want to use Hyper V or WSL, choose WSL.
    d. If your operating system is Windows, run Docker Desktop when it finished installing. Then, go to Settings->Resources-WSL Integration, and enable integration for WSL distro, in our case it's Ubuntu as shown below:
       ![Enable WSL Integration](image-3.png)
    e. Restart your computer

3. **Install Node.js with NPM**
    a. Please go to: https://nodejs.org/en and download the Node.js, recommend Node.js (LTS) v22.11.0 or any later LTS version
    b. Install the Node.js on your computer
    c. Open a terminal shell on your computer and verify that NPM is installed, run this command: npm -v

4. **Install Visual Studio Code IDE:**
    a. If you would like to develop with this project continue, then continue, if you just want to run the project using docker, proceed to step 4 (Install Git)
    b. You can use any IDE of your choice, but if you would like to use VS Code, download it here: https://code.visualstudio.com/
    b. Install the VS Code IDE

5. **Install Git:**
    a. Please download the Git package for your operating system using the link in the "Requirements" section
    b. Install Git

6. **Clone the MVP repository:**
    a. Navigate to the location where you would like to clone the project to
    b. Open git bash
    c. Run the following command: `git clone https://github.com/hunghpham/MVP-Application-Sensor-Telemetry-System.git`
    d. A folder of the project named "MVP-Application-Sensor_Telemetry-System" is created with the project cloned from github


## Usage

1. **Run The Project With Docker Only**
    a. Please run the Docker Desktop application
    b. If the OS is Windows, open Command Prompt or Windows power shell. If using other OS, use their terminal shell application
    c. Navigate in to where the folder of the project "MVP-Application-Sensor_Telemetry-System", the root of the project
    d. Run the docker command: `docker compose run --entrypoint="" frontend npm install`
    e. Wait until the installation of node modules are complete, it can take up to 5 minutes, since node modules is quite large
    f. Run the docker command: `docker compose up --build`
    g. Open a browser and enter the url: http://localhost:3000
    h. As the docker images are being built, keep waiting and the application will start up with simulation sensor data flowing through, you can hit the refresh button on the browser if your browser is caching something. 
       ![MVP Dashboard](image-3.png) 
    i. Done, you can now use the MVP application

2. **Run The Project For Development**
    a. Please run the Docker Desktop application
    b. Open an IDE of your choice, we will Visual Studio Code as example here
    c. In VS Code, File->Open Folder, and choose the "MVP-Application-Sensor_Telemetry-System" folder, which is the root of the project
    d. Open a terminal in VS Code
    e. Run the docker command: `docker compose run --entrypoint="" frontend npm install`
    f. Wait until the installation of node modules are complete, it can take up to 5 minutes, since node modules is quite large
    g. Run the dokcer command: `docker compose up --build`
    h. Open a browser and enter the url: http://localhost:3000
    i. As the docker images are being built, keep waiting and the application will start up with simulation sensor data flowing through, you can hit the refresh button on the browser if your browser is caching something.
       ![MVP Dashboard](image-3.png) 
    j. Done, you can now start developing with the MVP project

    

