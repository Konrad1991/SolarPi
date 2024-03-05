library(ggplot2)
set.seed(42)
ggplot() +
  geom_point(aes(x = sin(seq(0, 2 * pi, length.out = 100)),
                 y = cos(seq(0, 2 * pi, length.out = 100)),
                 size = 20),
             color = "lightgreen") +
  geom_text(aes(x = 0, y = 0, label = "SolarPi"), size = 20, fontface = "bold",
                vjust = 0.5) +
  theme_void() +
  theme(legend.position = "none")
ggsave("SolarPi_Logo.png", width = 6, height = 6, dpi = 300)

