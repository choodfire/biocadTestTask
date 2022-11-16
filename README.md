# Тестовое задание для стажировки в Biocad

Задание:
* Брать конфиги(host,port,password и т.д.) для соединения с БД (бд на выбор) а так же адрес директории
* Периодически осматривать директорию на наличие новых не обработанных еще файлов (.tsv) (вероятно держать в базе список уже обработанных) (см. лист "Исходные данные")
* Обработку файлов ставить в очередь
* Обработка файла - нужно распарсить .tsv и положить в соответствующую структуру (формат файла статичный, поля/количество не меняется)
* Данные из файла поместить в БД
* После обработки файла и записи в БД нужно сформировать файл(rtf,doc,pdf на выбор) с названием из поля *unit_guid* в входном файле, с данными по этому *unit_guid*
* Ошибки парсинга (например не соответсвие файла) - тоже записывать в БД и файл
* Выходные файлы размещать в отдельной директории
* Сделать API-интерфейс который позволит получать из БД данные с пагинацией (page/limit) для получения данных по *unit_guid*

---

## API-Интерфейс

Интерфейс возвращает JSON выбранных логов. Можно получить по команде `curl localhost:8080/data/{unit_guid}/{page}`