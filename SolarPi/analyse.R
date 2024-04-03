library(plotly)
setwd("/home/konrad/Documents/GitHub/SolarPi/SolarPi")
f <- "btsnoop_hci.log"
f <- "btsnoop_hci.log.old"
f <- "btsnoop_hci-3.log"
f <- "btsnoop_hci-4.log"
f <- "btsnoop_hci-nocon.log"
f <- "btsnoop_hci.log.last"
d <- readBin(f, what = "double", size = 4, n = file.info(f)$size,
             endian = "little") # 4, double
df <- data.frame(x = 1:length(d), y = d)
plot_ly(data = df[df$y < 0.5 & df$y > -0.5, ], x = ~ x, y = ~ y)

plot_ly(data = df[df$y < 20 & df$y > -5, ], x = ~ x, y = ~ y)


