English | [简体中文](./README_zh.md) | [Home](http://edgelog.devautoops.com) | [Documentation](http://edgelog.devautoops.com/help/)

# EdgeLog-Lightweight ELK

EdgeLog is a lightweight log management system, which is based on DevOps idea and is totally developed by ourselves.

Enterprise may pruduce many log statics when its application program is running, and the statics are very important and may grow rapidly. Managing the log statics is a essential part of everyday workflow. Many enterprises select ELK(Elasticsearch, Logstash, Kibana) suit to manage their log statics. ELK which is powerful can meet the needs of enterprises for log management, but there are also many inconveniences. For example, ELK components occupy a large amount of disk space,and it depends on various preset environments, also its data query efficiency is limited by the scalability of Elasticsearch, etc. These inconveniences make ELK more bulky. EdgeLog can overcome many inconveniences of ELK. EdgeLog is a lighter weight option for ELK, which can provide convenience for developers and enterprises in need.

Edgelog is developed in Golang, and It is mainly composed of the client application log collection component Agent and the server management component. The Agent'acquisition data is stored in a dedicated time series database. EdgeLog has excellent features based on Go language and time series database. EdgeLog has excellent performance, extremely convenient use ability, and extremely simple maintainability, and it can be used on different system platforms. For example, X86 and arm etc., and it also provide new ideas for massive data processing in the field of DevOps.

EdgeLog is lightweight and flexible in configuration and does not rely on any environmental restrictions. It has many excellent performance based on individual components. Query based on high-performance and scalable TDengine database which can greatly improve query efficiency. It has been verified that the query time of tens of millions of data is in milliseconds.

EdgeLog supports embedded environments well, which provides a solid foundation for the expansion of embedded hardware device capabilities. It also provides convenient reference, tools, and foundation for medium and large-scale operation management software to transition to hardware as well as enhancing product value.

## EdgeLog VS ELK
| Comparison item | EdgeLog | ELK |
| :---: | :---: | :---: |
| System occupancy| The main program is developed in Go language, and the storage database uses TDengine. The entire installation package is only tens of M in size. The installation and operation do not depend on other environments, and the client Agent's occupation of system resources can be ignored|The main program is developed in Java language, the installation package is relatively large, and the installation and operation are heavily dependent on the external environment, and the installation package is large and occupies a large amount of system resources|
| Retrieval performance|Efficient retrieval performance, cache-based and scalable time series database storage|Efficient retrieval performance|
| Storage performance|As long as the performance of the node machine is satisfactory, there will be no IO loss, and the Agent will fetch logs to the time series database in real time, with basically no delay |Reading IO files may result in log loss, and Logstash synchronization of logs may cause delays|
| Fully functional |A more lightweight log management system that can collect and manage various application logs|Heavyweight log management system|
| Embedded Systems|Support the integration of embedded systems, flexible and customizable|A small embeddable JavaScript engine|
| Scalability |Support cluster expansion|Support cluster expansion|
| Fault tolerance|The system has done further encapsulation and processing of log data considering the needs of the enterprise|The data mining ability is weak. If you need to meet the data requirements of the enterprise, you need to do in-depth development of ELK|
| Front-end operation|There is a beautiful operation front-end, the operation is simple and clear, and you only need to click the mouse to complete the addition of related nodes, log retrieval and various aggregation queries|Front-end operations rely heavily on kibana performance and can complete search and aggregation functions|


## EdgeLog features

- [Lightweight] -EdgeLog is a compact and functional log management system. It is adopted by C/S mode and B/S mixed mode. It has the characteristics of excellent writing speed, high storage capacity, low resource occupation, and low deployment and maintenance costs. The purpose of its birth is to solve the problem of avoiding excessive use and operation and maintenance costs for the log system in the enterprise business. It does the most essential thing of a logging system with a simple design. It is aiming at efficient and stable operation of the system, low cost of use and service operation and maintenance costs. The sections of the system can be scaled horizontally to ensure high availability, high performance, and high capacity, stable functional design, controllable network bandwidth usage, and low memory usage.

- [Modular] -EdgeLog consists of client-side Agent and server-side main program, and both components are written in Go language. Nginx acts as the management background web server. The configuration of each component is flexible and simple, occupying less system resources.

- [More efficiency query] -EdgeLog uses TDengine to store logs collected in real time. The query is based on the high-performance and scalable TDengine. It uses optimized and efficient query sql to perform various aggregation calculations, and the computing resources are less than 1/5 of the general big data solution.

- [Mass storage] -EdgeLog uses a dedicated time series database TDengine to store massive log data. TDengine chooses columnar storage and advanced compression algorithm, the storage space is less than 1/10 of the general database, the storage data is stored for one year by default, and users can configure the storage time according to their needs, and based on TDengine's efficient cluster expansion scheme to ensure high data availability At the same time, it can also release the pressure on the back-end storage system to meet the specific needs of enterprises.

- [Easy to user] -EdgeLog provides a web-side system management portal. The background menu layout is simple and clear. Each menu has independent functions and simple configuration. The system cooperates with the central management screen to display various application data in aggregated graphics, which are flexible and changeable.

- [Multi-platforms, Compatible with mainstream CPUs and architectures] -Support Linux, Windows and other mainstream operating system platforms. Today, when vigorously advocating information technology application innovation and promoting the localization of information security, EdgeLog can also well support mainstream OS architectures, supporting X86, arm and other architectures. The designed Agent module has independent executable file and configuration file. Opening the corresponding collection module in the configuration file can realize the collection of enterprise application logs and truly realize the function of flexible configuration.

- [Wide range of applicable scenarios and scalable application] -EdgeLog has its own unparalleled advantages in log collection and management. For logs of various enterprise applications (such as Nginx, Apache, Tomcat, MySQL, etc.), EdgeLog can flexibly select the log management module unit of the corresponding system according to different applications, and EdgeLog provides different interfaces for logs generated by different applications. How to manage these logs, EdgeLog flexibly calls the corresponding hive interface according to the ideas of enterprise administrators, which is convenient and time-saving, and greatly improves the production efficiency of enterprises. The log collection adopts the C/S service mode, the management background adopts the B/S mode, and the Agent log collection client is suitable for various operating system platforms. Through the EdgeLog system, you can customize the log indicators that customers need to pay attention to, support the search for keyword log information from a large number of logs, and timely feedback the workload and performance of the application in the customer's production system, which can timely warn of various production safety hazards in the system. And according to the needs of the enterprise, the relevant functional modules can be customized for the special application of the enterprise to achieve flexibility and scalability.

To learn more about the EdgeLog management system, you can visit the official website of [EdgeLog](http://edgelog.devautoops.com).

## Standard Self-Check for Excellent DevOps Tools
| Characteristic | EdgeLog |
| --- | --- |
| Easy to integrate|It is easy to integrate with other components and can also be implemented well in embedded systems|
| Powerful API support | Rich API and detailed documentation|
| Support multi platforms| Support Linux, Windows and other mainstream operating system platforms, and also widely support common chip platforms |
| Software development automation|Supports iterative updates and automatic distribution of software components|
| Customizable|Provide open source version and professional version, and provide professional customized services for specific scenarios|
| Simple to use|The system function modules are simple and easy to use, and each component has a detailed help guide and rich API.|
| Dashboard|Beautiful dashboard with common parameters by default|
| High performance|Both the client and the server are developed in Go language, and the data storage is based on a scalable and efficient query time series database. The system combines the advantages of both, and the overall performance is extremely powerful.|
| Price|Free to use open source version|
| Support CI/CD|Based on DevOps concept development and continuous iteration, support continuous integration and continuous distribution|
| Customer Support|Professional team, rich experience in telecommunications and financial high-standard enterprise services, quick response|

## Demo address

Experience EdgeLog Open Source Edition Online [Demo System](http://demo.edgelog.devautoops.com)。

## System screenshot
![Dashboard](./dashboard.png)

![Log_analysis](./log_analysis.png)

![Node_monitor](./node_monitor.png)

## Open Source Instructions

The source code of the server has been open sourced in the current Git library, and detailed API documentation is provided. For the source code of the client Agent, please refer to [Agent Open Source Library](https://github.com/). If you have any questions, you can submit issues and we will reply as soon as possible.
If you want to install and try the complete system locally, you can [download](http://edgelog.devautoops.com/download/) the complete front-end and back-end installers.

## API documentation

To understand the detailed APIs of the system, you can visit [EdgeLog API](http://edgelog.devautoops.com/help/api.html)。

## Installation & Operation & Help

Please follow the [installation documentation](http://edgelog.devautoops.com/help/install.html) to install EdgeLog, For operation manuals and help documents, please read [Documentation & Help](http://edgelog.devautoops.com/help/).

## Professional help and support

We continue to iterate according to the needs of the market and enterprises, and are also willing to provide professional customized services based on the particularity of the scene。If necessary, go to [Contact Us](http://edgelog.devautoops.com).

## Others

Agent is another part of EdgeLog system, and it is designed to collect log statics from application program running on client machine. If you want to konw more about Agent, [click](https://github.com/DevAutoOps/EdgeLog-Agent).


