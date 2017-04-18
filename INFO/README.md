### Telegraf
Для сбора различной метрики в проекте используется open-source сервис телеграф https://github.com/influxdata/telegraf 
Сервис построен на плагинах. Существуют input и output плагины. 
Input плагины позволяют получать определенную метрику системы (состояние оперативной памяти, нагрузка на ядра процессора, и т.д.). Output плагины 
релизовывают вывод собранной метрики куда-либо (файл, influxdb, и т.д.). Для нашего проекта был написан плагин вывода в mongodb (https://github.com/grabgo/telegraf/tree/mongo-output-plugin-custom)
Проект собирается с помощью менеджера зависимостей gdm (go get github.com/sparrc/gdm) командой make в бинарный файл. 
Внимание сборка проекта может вызвать затруднения, если скачать его с github.com/influxdb/telegraf, так как в этом репозитории не будет необходимых зависимостей.
Вместо этого лучше скачать проект с github.com/influxdata/telegraf и связать директорию telegraf с https://github.com/grabgo/telegraf/ для того, чтобы собрать бинарник с плагином вывода в mongodb.

После установки необходимо правильно сконфигурировать сервис. Это делается с помощью конфигурационного файла (https://github.com/influxdata/telegraf/blob/master/docs/CONFIGURATION.md)
По умолчанию файл лежит в директории /etc/telegraf/telegraf.conf. Он поделен на разделы, нас интересует ввод и вывод. 
 
Настойка метрики, которую мы будем выводить выглядит так:   
###############################################################################
#                            INPUT PLUGINS                                    #
###############################################################################

# Read metrics about cpu usage
[[inputs.cpu]]
  ## Whether to report per-cpu stats or not
  percpu = true
  ## Whether to report total system cpu stats or not
  totalcpu = true
  ## If true, collect raw CPU time metrics.
  collect_cpu_time = false


# Read metrics about memory usage
[[inputs.mem]]
  # no configuration

# Get the number of processes and group them by status
[[inputs.processes]]
  # no configuration
  
# Read metrics about system load & uptime
[[inputs.system]]
  # no configuration
  
  
Настройка вывода в mongodb выглядит так:
###############################################################################
#                            OUTPUT PLUGINS                                   #
###############################################################################

## MongoDB host. If port is not specified, then default (27017) will be used for a connection.
hosts = ["localhost:27017"]    #required (at least one)
## MongoDB database name   #required
db = "mongodb_name"   #required
## MongoDB collection name, where documents are or will be stored.
collection = "telegraf_metric"  #required
## User credentials. Not required.
username = "username"
password = "password"
## Just a unique field for environment in case you want to send similar by type data to one mongodb from different servers
server_name = "server_name"
  
В hosts через запятую необходимо указать хосты mongodb (по крайней мере один), параметр обязателен
db - имя базы данных в mongo, параметр обязателен
collection - обязательный параметр, имя коллекции, должно совпадать с именем коллекции указанной в properties файле, используемом в MongoPersistenceConfig
username и password - необязательные параметры авторизации
server_name - этот параметр должен быть различным на каждом из серверов, именно с помощью этого параметра будут различаться метрики разных серверов одного пула.
Рекомендуется избегать значений "server" и "metric" для этого параметра, чтобы не возникало двусмысленных ситуаций в rest api при получении метрик. 
Все сервера должны иметь свободный доступ к mongodb. Данный плагин не осуществляет подключение по ssh. 
После запуска плагина необходимо проверить запись в монго. В админке во вкладке telegraf: проверить кол-во серверов и определенные на каждом из них метрики.
Преимущество в том, что на каждом из серверов можно сконфигурировать свои метрики и все они будут разделены для чтения из mongodb. 

###Установка
В репозитории лежит deb пакет версии 1.2.1. Установить его и заменить бинарный файл (обычно устанавливается по пути /usr/bin) на собраный бинарник с плагином.
service restart telegraf
Если запустить бинарный файл вручную (не в качестве сервиса), то вывод будет следующим
2017/04/18 12:10:16 I! Using config file: /etc/telegraf/telegraf.conf
2017-04-18T09:10:16Z I! Starting Telegraf (version 1.2.0-rc1-127-g516dffa)
2017-04-18T09:10:16Z I! Loaded outputs: mongodb
2017-04-18T09:10:16Z I! Loaded inputs: inputs.cpu inputs.mem inputs.processes inputs.system
2017-04-18T09:10:16Z I! Tags enabled: host=yuri-desktop
2017-04-18T09:10:16Z I! Agent Config: Interval:10s, Quiet:false, Hostname:"yuri-desktop", Flush Interval:10s 
^C2017-04-18T09:10:31Z I! Hang on, flushing any cached metrics before shutdown