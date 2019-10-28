ipb4parser v.0.0.4
------------------

Merge-парсер для форумов на движке IPB4. Объединяет тему или посты пользователя в один удобный файл, игнорируя сообщения без лайков, для последующей обработки.

example:

ipb4parser.exe --topic [url]

ipb4parser.exe --topic https://forum.bits.media/index.php?/topic/5018-%D0%B1%D0%B8%D1%82%D0%BA%D0%BE%D0%B8%D0%BD%D1%8B-%D0%B2%D0%B0%D0%BB%D1%8E%D1%82%D0%B0-%D0%B1%D1%83%D0%B4%D1%83%D1%89%D0%B5%D0%B3%D0%BE/ --count 2
ipb4parser.exe --user https://forum.bits.media/index.php?/profile/29870-chingizzz/

flags:

--topic 
--user
--count | default=0, i.e > 0

Issues:

Компиляция не работает из-за некорректно настроенного GOPATH во время разработки на прошлой ОС. 
