#include <LiquidCrystal_PCF8574.h>
#include <ArduinoOTA.h>
#include <PubSubClient.h>
#include <Wire.h>
#include "DEVICE.h"
#include "MONITOR.h"


LiquidCrystal_PCF8574 lcd(0x3F); // set the LCD address to 0x27 for a 16 chars and 2 line display


WiFiManager   wm;                 // global wm instance
//WiFiClient    espClient;
WiFiClientSecure espClient;
PubSubClient  client(espClient);

DEVICE        device;
MONITOR       monitor;


char mqtt_server[]      = "w.wiin.win";
int mqtt_port           = 1884;
char mqtt_user[]        = "mao";
char mqtt_password[]    = "linmao8888";

int show = -1;

const int PIN_BEEP = 15;
bool beep_on = false;

void setup()
{
  Serial.begin(115200, SERIAL_8N1);
  Serial1.begin(115200, SERIAL_8N1); //   RX-GPIO3-? TX-GPIO4-D4
  IPRINTLN("\n\nBootup...\n");


  String mac = WiFi.macAddress();
  device.setup(mac.c_str());

  setup_wifi(wm);


  Wire.begin();
  Wire.beginTransmission(0x3F);
  int error = Wire.endTransmission();
  if (error == 0) {
    IPRINTLN(" - LCD found.");
    lcd.begin(16, 2); // initialize the lcd
  } else {
    IPRINTLN(" - LCD not found.");
  }


  lcd.setBacklight(255);
  lcd.home();
  lcd.clear();
  lcd.setCursor(0, 0);
  lcd.print("Booting...");
  
  pinMode(PIN_BEEP, OUTPUT);
  digitalWrite(PIN_BEEP, LOW);
  IPRINTLN("Device Ready...");

} // setup()

static const char *fingerprint PROGMEM = "44 14 9A 3F C3 E9 F1 F3 84 1A B4 9F B6 4D 19 8A B2 92 31 D6";

void loop()
{
  static int last_wifi_status = WL_DISCONNECTED;
  //DPRINT(".");
  //DPRINT("FREE HEAP:%d\n",ESP.getFreeHeap());

  checkButton(wm);
  int current_wifi_status = WiFi.status();

  if (current_wifi_status == WL_CONNECTED) {                  //WIFI 已连接
    if (last_wifi_status != WL_CONNECTED) {                   //WIFI 刚连接上
      led_off();
      //strcpy(device_ip, WiFi.localIP().toString().c_str());
      device.setIP(WiFi.localIP().toString().c_str());
      IPRINTF(" done, IP: %s\n", device_ip);
      setup_OTA();
    }

    ArduinoOTA.handle();
    if (client.connected()) {               // mqqt client 连接成功
      DPRINTF(".");
      led_off();
      dojob();
      client.loop();
    } else {                                // mqqt client 连接断开
      DPRINTF("-");
      show_offline("MQTT connecting ");
      led_on();
      IPRINTLN();
      //IPRINTF("MQTT client offline\n");
      //IPRINTF("MQTT client connecting to %s:%d\n", mqtt_server, mqtt_port);

      //espClient.setFingerprint(fingerprint);
      espClient.setInsecure();
      client.setBufferSize(512);
      client.setKeepAlive (30);
      client.setSocketTimeout (30);
      client.setServer(mqtt_server, mqtt_port);
      client.setCallback(mqttcallback);
      if (client.connect(device_id_full, mqtt_user, mqtt_password, topic_will, 2, true, "offline")) {

        bool isOK = publishMQTTMessage( client, topic_will, "online", true);
        if (!isOK) {
          client.disconnect();
        }
        IPRINTLN("connected");
        IPRINTF("publish topic_will :\n   - %s\n", topic_will);
        IPRINTLN("Subscribe Topic : ");

        //meter.subscribe_topics(client);
        monitor.mqtt_subscribe_topics(client);
        device.mqtt_subscribe_topics(client);

        IPRINTLN();
      }
#ifdef DEBUG
      else {
        IPRINT("failed, rc=");
        IPRINTLN(client.state());

      }
#endif
    }
  }
  else {                                    //  WIFI 断开
    led_on();
    //DPRINTF("=");
    //delay(750);
    WiFi.disconnect();
    WiFi.hostname(device.getDeviceID());
    //WiFi.setPhyMode(WIFI_PHY_MODE_11B);
    WiFi.mode(WIFI_STA);
    //WiFi.begin(wifi_ssid, wifi_password);
    WiFi.begin();
    //delay(750);

    int count = 40;
    IPRINT("\nWiFi connecting to ");
    show_offline("WiFi connecting ");
    IPRINT(WiFi.SSID());
    while (WiFi.status() != WL_CONNECTED && (count-- > 0)) {  // 等待4x500=2000ms
      IPRINT("=");
      checkButton(wm);
      //checkButton(device_id_hostname);
      delay(500);
    }
  }


  //  long loop_cycle_millis = 1000 * 10; //5s
  //  static long last_loop_finished_millis = 0;
  //
  //  long current_loop_finish_millis = millis();
  //  long delay_millis =  last_loop_finished_millis + loop_cycle_millis - current_loop_finish_millis;
  //  if (delay_millis < 0) delay_millis = 0;
  //  if (delay_millis > loop_cycle_millis)  delay_millis = loop_cycle_millis;
  //  last_loop_finished_millis = current_loop_finish_millis;
  //DPRINTF("delay_millis = %d\n",delay_millis);

  last_wifi_status = current_wifi_status;
  //ESP.wdtFeed();
  delay(250);
}

void show_offline(char* msg) {
  lcd.setBacklight(50);
  lcd.home();
  lcd.clear();
  lcd.setCursor(0, 0);
  lcd.print(msg);
}

void dojob() {

  if (monitor.invalid) {
    int t_max = 0;
    for(int i=0;i<3;i++){
       if(monitor.t[i]>t_max) t_max = monitor.t[i];
    }
    sprintf(payload, "W %d   %04dW  %02C", monitor.workers,monitor.ptotal,t_max);
    //    lcd.setBacklight(50);
    //    lcd.home();
    //    lcd.clear();
    lcd.setCursor(0, 0);
    lcd.print(payload);
    //lcd.blink();

    sprintf(payload, "%.9f|%.3f",  monitor.unpaid, monitor.forcast_24h);
    lcd.setCursor(0, 1);
    lcd.print(payload);

    if(monitor.workers<3 || t_max >75) beep_on = true;
    else beep_on=false;

    monitor.invalid = false;
  }
  int beep_loop = 4;
  static int s = 0;
  if (s == 0) {
    lcd.setCursor(1, 0);
    lcd.print(":");
    if (beep_on) {
      digitalWrite(PIN_BEEP, HIGH);
    }
  }
  if (s == beep_loop/2) {
    lcd.setCursor(1, 0);
    lcd.print(" ");
    digitalWrite(PIN_BEEP, LOW);
  }
  if (s++ == beep_loop) s = 0;


  
  static volatile uint32_t lastMillis                  = 0;
  static volatile uint32_t lastConfigMessageTimeStamp  = 0;
  static volatile uint32_t lastStateMessageTimeStamp   = 0;

  long now = millis();
  if (now - lastStateMessageTimeStamp > 1000 * device.cycle_seconds) { //
    //IPRINTF("gap = %d\n",now - lastStateMessageTimeStamp);
    lastStateMessageTimeStamp = now;
    device.pub_signal = device.pub_signal | PUB_STATE_ALL;

    if (now - lastConfigMessageTimeStamp > 1000 * 60 * 5) { //
      lastConfigMessageTimeStamp = now;
      device.pub_signal = device.pub_signal | PUB_CONFIG;
    }
  }

  if (device.pub_signal & PUB_CONFIG) {
    //DPRINTF("meter.publish_config\n");
    //monitor.publish_config(client);
    //DPRINTF("device.publish_config\n");
    device.publish_config(client);
    device.pub_signal = device.pub_signal ^ PUB_CONFIG;
  }
  if (device.pub_signal & PUB_STATE_DEVICE) {
    //DPRINTF("device.publish_state\n");
    device.publish_state( client,WiFi.SSID().c_str(), WiFi.RSSI());
    device.pub_signal = device.pub_signal ^ PUB_STATE_DEVICE;
  }
//  if (device.pub_signal & PUB_STATE_SENSORS) {
//    device.publish_state( client,WiFi.SSID().c_str(), WiFi.RSSI());
//    device.pub_signal = device.pub_signal ^ PUB_STATE_SENSORS;
//  }
}


void mqttcallback(char* intopic, byte* in_payload, unsigned int length) {
  led_on();
  IPRINTLN();
  IPRINTF("received - %s\n", intopic);

  device.mqtt_callback(intopic, in_payload, length);
  monitor.mqtt_callback(intopic, in_payload, length);

  led_off();
}


void display() {
  Serial.println(show);
  if (show == 0) {
    lcd.setBacklight(255);
    lcd.home();
    lcd.clear();
    lcd.print("Hello LCD");
    delay(1000);

    lcd.setBacklight(0);
    delay(400);
    lcd.setBacklight(255);

  } else if (show == 1) {
    lcd.clear();
    lcd.print("Cursor On");
    lcd.cursor();

  } else if (show == 2) {
    lcd.clear();
    lcd.print("Cursor Blink");
    lcd.blink();

  } else if (show == 3) {
    lcd.clear();
    lcd.print("Cursor OFF");
    lcd.noBlink();
    lcd.noCursor();

  } else if (show == 4) {
    lcd.clear();
    lcd.print("Display Off");
    lcd.noDisplay();

  } else if (show == 5) {
    lcd.clear();
    lcd.print("Display On");
    lcd.display();

  } else if (show == 7) {
    lcd.clear();
    lcd.setCursor(0, 0);
    lcd.print("*** first line.");
    lcd.setCursor(0, 1);
    lcd.print("*** second line.");

  } else if (show == 8) {
    lcd.scrollDisplayLeft();
  } else if (show == 9) {
    lcd.scrollDisplayLeft();
  } else if (show == 10) {
    lcd.scrollDisplayLeft();
  } else if (show == 11) {
    lcd.scrollDisplayRight();

  } else if (show == 12) {
    lcd.clear();
    lcd.print("write-");

  } else if (show > 12) {
    lcd.print(show - 13);
  } // if

  delay(1400);
  show = (show + 1) % 16;
} // loop()
