package config

import "time"

// TimeInterval - временной интервал для контроля частоты запросов
var TimeInterval time.Duration = time.Second * 1

// CacheCleanerDuration - интервал времени для очистки устаревших записей в кэше
var CacheCleanerDuration time.Duration = time.Second * 3

// MaxRequestsPerTimeInterval - максимальное количество запросов в заданный временной интервал
const MaxRequestsPerTimeInterval int = 7

// WriteLogs - флаг указывающий, нужно ли вести логирование
const WriteLogs bool = true

// CleanCache - флаг указывающий, нужно ли производить очистку кэша
const CleanCache bool = true
