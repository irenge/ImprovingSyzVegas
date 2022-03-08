hours  <- c(0, 0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5, 5.5)
                #6, 6.5)
coverage <- c( 0, 20867, 32094,40146, 45361, 49074, 52384, 57919, 61656, 63724, 66072, 68748)
#, 70751, 72162)
executed <- c(4, 16887, 34177, 50857, 67793, 84879, 101877, 118553, 135152, 151270,166441, 179638)
#,196036, 212761)
crashes <- c(0, 0, 0, 0, 2, 7, 9, 12, 16, 18, 23, 36)
                 #, 36, 42)
plot(hours, coverage, type = "o", lty=1, col="red")
#lines(coverage, type = "o", col = "blue")
#### 
execut <- c(4, 16887, 34177, 50857, 67793, 84879, 101877, 118553, 135152, 151270,166441, 179638)

six_hours <- data.frame(hours, coverage, executed, crashes)

ggplot(six_hours, aes(x=hours, y=coverage)) + geom_line(aes(color=crashes), size=4 ) + geom_path()+ geom_point(color="red", size=3) 

## 5 hours with Vegas

hours <- c(0, 0.5, 1, 1.5, 2, 2.5, 3,  3.5, 4, 4.5, 5, 5.5)
coverage<-c(0,13672,22106,  28964, 32742, 35639,38281, 43046, 45393, 47777,50038, 51204)
executed <-c(5,16868,33784,50526, 67505, 83871, 100116,116678, 133275, 150156, 166598, 180327)
crashes <- c(0, 0, 2, 5, 6, 13, 20, 24, 28, 34, 39, 40)
               

plot(hours, coverage, type = "o", lty=1, col="blue")

six_hours_Vegas <- data.frame(hours, coverage, executed, crashes)
ggplot(six_hours_Vegas, aes(x=hours, y=coverage)) + geom_line(aes(color=crashes), size=4 ) + geom_path()+ geom_point(color="red", size=3) 

## combined 

hours <- c(0, 0.5, 1, 1.5, 2, 2.5, 3,  3.5, 4, 4.5, 5, 5.5)
executed_v <-c(5,16868,33784,50526, 67505, 83871, 100116,116678, 133275, 150156, 166598, 180327)
executed <- c(4, 16887, 34177, 50857, 67793, 84879, 101877, 118553, 135152, 151270,166441, 179638)
crashes_v <- c(0, 0, 2, 5, 6, 13, 20, 24, 28, 34, 39, 40)
coverage_v<-c(0,13672,22106,  28964, 32742, 35639,38281, 43046, 45393, 47777,50038, 51204)

coverage <- c( 0, 20867, 32094,40146, 45361, 49074, 52384, 57919, 61656, 63724, 66072, 68748)
executed <- c(4, 16887, 34177, 50857, 67793, 84879, 101877, 118553, 135152, 151270,166441, 179638)
crashes <- c(0, 0, 0, 0, 2, 7, 9, 12, 16, 18, 23, 36)
combined <- data.frame(hours, executed_v, executed, coverage, coverage_v, crashes, crashes_v)

ggplot(combined, aes(x=hours)) + 
  geom_line(aes(y = executed_v), color = "darkred", size=4) + 
  geom_line(aes(y = executed), color="steelblue", linetype="twodash", size=4) 

ggplot(combined, aes(x=hours)) + 
  
  geom_line(aes(y = coverage), color = "darkred", size=4) + 
  geom_line(aes(y = coverage_v), color="steelblue", linetype="twodash", size=4)+
  coord_flip(ylim = c(0,70000), xlim = c(0, 4))


ggplot(combined, aes(y=hours)) + 
  
  geom_line(aes(x = coverage), color = "darkred", size=4) + 
  geom_line(aes(x = coverage_v), color="steelblue", linetype="twodash", size=4)+
  coord_flip(xlim = c(3,40000), ylim = c(0, 4))

ggplot(combined, aes(y=executed)) + 
  ggtitle("Coverage per execution: Syzkaller Vs SyzVegas") +
  xlab("Coverage") + ylab("Execution") +
  geom_line(aes(x = coverage), color = "darkred", size=4) + 
  geom_line(aes(x = coverage_v), color="steelblue", linetype="twodash", size=4) +
#+
 # geom_label()
 
 #scale_x_discrete(name ="Coverage") +
 # scale_y_discrete(name ="Execution") +
  theme(axis.text.x = element_text(face="bold", color="#993333", 
                                   size=14, angle=45),
        axis.text.y = element_text(face="bold", color="#993333", 
                                   size=14, angle=45)) +
  theme( axis.line = element_line(colour = "darkblue", 
                                  size = 1, linetype = "solid")) + 
  geom_label(label="Syzkaller",
             x =40000,
             y= 50000) +
  geom_label(label="SyzVegas",
             x =40000,
             y= 100000)
 # coord_flip(xlim = c(3,40000), ylim = c(0, 4))


# 24 hr with Vegas

hours <- c(0, 0.5, 1, 2, 3, 4, 5 , 6 ,7, 8, 9, 10, 11, 12 ,13, 14, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28 )
coverage <- c(0, 14228, 29450, 36153,39681, 40118, 44068, 47327,51000, 54359,56143, 59710, 60905, 62023, 63539, 65063, 66358,  67984, 68697, 69826,70318,70981,71477, 72021,72268, 72779,73319, 74129,  74341 )
executed <- c(3,17883, 53228, 85688, 116627, 120188, 151381, 181584,210687, 239380,269274, 300281,330627, 358834, 389779, 422172,  451068,  482743, 512683, 543944, 574633,  598153, 625749, 652759, 682924,712136, 745800,803804, 833963)
crashes <- c(0, 0, 3,8, 11, 17, 29, 35, 40, 54, 58, 76, 86, 96, 103, 124, 138, 162,175,184,196, 205, 237, 252, 276, 287, 316, 331, 345)

twenty_four_hours_Vegas <- data.frame(hours, coverage, executed, crashes)

plot(syz_24hours, syzV_24cov, type = "o", lty=1, col="green")

library(ggplot2)
ggplot(twenty_four_hours_Vegas, aes(x=hours, y=coverage)) + geom_point() + geom_smooth(method="lm")

ggplot(twenty_four_hours_Vegas, aes(x=hours, y=coverage)) + geom_line(aes(color=crashes), size=4 ) + geom_path()+ geom_point(color="red", size=3) 




ggplot(combined, aes(x=executed)) + 
  ggtitle("Coverage per execution: Syzkaller Vs SyzVegas") +
  ylab("Coverage") + xlab("Execution") +
  geom_line(aes(y = coverage), color = "darkred", size=4) + 
  geom_line(aes(y = coverage_v), color="steelblue", linetype="twodash", size=4) +
  #+
  # geom_label()
  
  #scale_x_discrete(name ="Coverage") +
  # scale_y_discrete(name ="Execution") +
  theme(axis.text.x = element_text(face="bold", color="#993333", 
                                   size=14, angle=45),
        axis.text.y = element_text(face="bold", color="#993333", 
                                   size=14, angle=45)) +
  theme( axis.line = element_line(colour = "darkblue", 
                                  size = 1, linetype = "solid")) + 
  geom_label(label="Syzkaller",
             y =40000,
             x= 50000) +
  geom_label(label="SyzVegas",
             y =40000,
             x= 100000)
# coord_flip(xlim = c(3,40000), ylim = c(0, 4))