
#include "DEVICE.h"
#include <ArduinoOTA.h>
#include <EEPROM.h>

extern DEVICE device;


const char sw_version[30]         = "2022.05.05_2.1";   //        = ;//
const char device_manufacture[20] = "Gaoyan";               //
const char device_class[20]       = "Monitor";

char device_id[20];         //    = "XXXXXXXXXXXXXXXXXX";
char device_id_full[36];    //  = "KP-PowerMeter-XXXXXXXXXXXXXXXXXXXXXX";
char device_name[36];       //  = "KP PowerMeter XXXXXXXXXXXXXXXXXXXXXX";
char device_model[36];      //  = "PowerMeter XXCH ESP8266";

char device_ip[20];         //  = "000.000.000.000";
char device_mac[20];        //  = "XX:XX:XX:XX:XX:XX";

char topic_will[50];        //   = "haworkshopyc1/sensor/XXXXXXXXXXXXXXXXXXXXXXXXXX/status";
char topic_state_device[50];//       = "haworkshopyc1/sensor/XXXXXXXXXXXXXXXXXXXXXXXXXX/state";

char topic_state[50];       //         = "haworkshopyc1/sensor/XXXXXXXXXXXXXXXXXXXX/XXXXXXXXXXXXXXXXXX/state";
char topic_config[50];      //        = "haworkshopyc1/switch/XXXXXXXXXXXXXXXXXXXX/XXXXXXXXXXXXXXXXXX/config";
char topic_set[50];         //           = "haworkshopyc1/switch/XXXXXXXXXXXXXXXXXXXX/XXXXXXXXXXXXXXXXXX/set";

//////////////////////////////////////////
char payload[1000];
char device_config_json[300];



void DEVICE::mqtt_subscribe_topics(PubSubClient &client) {
  sprintf(topic_set,    "haworkshopyc1/number/%s/%s/set",     "Cycle", device_id);
  client.subscribe(topic_set);
  IPRINTF("   - %s\n", topic_set);

  sprintf(topic_set,    "haworkshopyc1/button/%s/%s/set",     "Restart", device_id);
  client.subscribe(topic_set);
  IPRINTF("   - %s\n", topic_set);

  
}


bool DEVICE::mqtt_callback(char* intopic, byte* in_payload, unsigned int length) {  
  
  sprintf(topic_set,    "haworkshopyc1/button/%s/%s/set",     "Restart", device_id);
  if (!strcmp(intopic, topic_set)) {
    if (!strncmp("restart", (char*)in_payload, length)) {
      IPRINTF("received command - Restarting\n\n");
      delay(10);
      ESP.restart();
    }
  }
  sprintf(topic_set,    "haworkshopyc1/number/%s/%s/set",     "Cycle", device_id);
  if (!strcmp(intopic, topic_set)) {
    strncpy(payload, (char*)in_payload, length);
    payload[length] = '\0';
    int v = atoi(payload);
    if (v <= 300 && v >= 5) {
      device.cycle_seconds = v;
      //update_all_state = 1;//  update_switch_state
      device.pub_signal = device.pub_signal | PUB_STATE_DEVICE;
      IPRINTF("received command - Set Cycle Seconds = %d\n", device.cycle_seconds);
      byte vv[2] = {v >> 8, v & 0xFF};
      EEPROM_rotate_write(vv, 2);
    }
    return true;
  }

  return false;
}





void led_on() {
  //digitalWrite(LED_PIN, LOW);
}
void led_off() {
  //digitalWrite(LED_PIN, HIGH);
}


// 数据的当前地址存在510，511处
bool EEPROM_rotate_init() {
  EEPROM.begin(512);
  int eeprom_address = EEPROM.read(510) * 256 + EEPROM.read(511);

  if (eeprom_address < 0 || eeprom_address > 509) {
    eeprom_address = 0;
    EEPROM.write(510, (eeprom_address >> 8) & 0xFF);
    EEPROM.write(511, eeprom_address & 0xFF);

    if (EEPROM.commit()) {
      IPRINTF("committed address 0x%X to 510/511\n", eeprom_address);
    } else {
      IPRINTLN("commit FAILED to 510/511!! system cruppted");
      return false;
    }
  }
  return true;
}
bool EEPROM_rotate_read(byte* v, int length) {
  int eeprom_address = EEPROM.read(510) * 256 + EEPROM.read(511);

  if (eeprom_address < 0 || eeprom_address > 509) {
    eeprom_address = 0;
    EEPROM.write(510, (eeprom_address >> 8) & 0xFF);
    EEPROM.write(511, eeprom_address & 0xFF);

    if (EEPROM.commit()) {
      IPRINTF("committed address 0x%X to 510/511\n", eeprom_address);
    } else {
      IPRINTLN("commit to 510/511 failed!! system must be cruppted");
      return false;
    }
  }
  for (int i = 0; i < length; i++) {
    *(v + i) = EEPROM.read((eeprom_address + i > 509) ? (eeprom_address + i - 510) : (eeprom_address + i));
  }

  return true;
}
bool EEPROM_rotate_write(byte* v, int length) {

  static int eeprom_commit_count = 0;
  //byte high = EEPROM.read(510);
  //byte low = EEPROM.read(511);
  //int eeprom_address =  (int)high << 8 + (int)low;
  int eeprom_address = EEPROM.read(510) * 256 + EEPROM.read(511);

  //Serial.printf("eeprom_address:%x\n", eeprom_address);


  if (eeprom_address < 0 || eeprom_address > 510) {
    IPRINTF("System must be cruppted\n");
  }
  else {
    if (eeprom_commit_count > 509) {
      eeprom_address = eeprom_address + 1;
      if (eeprom_address > 509)       eeprom_address = 0;


      EEPROM.write(510, (eeprom_address >> 8) & 0xFF);
      EEPROM.write(511, eeprom_address & 0xFF);

      if (EEPROM.commit()) {
        IPRINTF("committed address 0x%X to 510/511\n", eeprom_address);
      } else {
        IPRINTLN("commit address failed to 510/511!! system cruppted");
      }


      eeprom_commit_count = 0;
    }

    //Serial.printf("EEPROM 510:%X,511:%X \n", EEPROM.read(510), EEPROM.read(511));

    Serial.printf("EEPROM(Addr 0x%X): ", eeprom_address);

    for (int i = 0; i < length; i++) {

      EEPROM.write( ((eeprom_address + i > 509) ? (eeprom_address + i - 510) : (eeprom_address + i)), *(v + i));

    }

    if (EEPROM.commit()) {
      eeprom_commit_count ++;
      IPRINTF("committed,count:%d\n", eeprom_commit_count);
    } else {
      IPRINTLN("commit FAILED!!");
      return false;
    }
  }

  return true;
}

void setup_OTA() {

  // Port defaults to 8266
  // ArduinoOTA.setPort(8266);

  //String mqtt_clientid = String(DEVICE_CLASS)+"-"+DEVICE_NAME;
  //sprintf(DEVICE_IDENTIFIER, "%s-%s", DEVICE_CLASS, DEVICE_NAME);
  // Hostname defaults to esp8266-[ChipID]
  ArduinoOTA.setHostname(device.getDeviceID());

  // No authentication by default
  // ArduinoOTA.setPassword("admin");

  // Password can be set with it's md5 value as well
  // MD5(admin) = 21232f297a57a5a743894a0e4a801fc3
  // ArduinoOTA.setPasswordHash("21232f297a57a5a743894a0e4a801fc3");
  IPRINTLN();

  ArduinoOTA.onStart([]() {
    if (ArduinoOTA.getCommand() == U_FLASH) {
      //type = "sketch";
      IPRINTLN("\n");
      IPRINTLN("======================================================================");
      IPRINTLN("Start updating sketch");
    } else { // U_FS
      //type = "filesystem";
      IPRINTLN("\n");
      IPRINTLN("======================================================================");
      IPRINTLN("Start updating filesystem");
    }

    // NOTE: if updating FS this would be the place to unmount FS using FS.end()
    //IPRINTLN("Start updating " + type);
  });
  ArduinoOTA.onEnd([]() {
    IPRINTLN("\nEnd");
  });
  ArduinoOTA.onProgress([](unsigned int progress, unsigned int total) {
    static long pt = 0;
    if (millis() - pt > 200 || progress > total - 10) {
      IPRINTF("%u%%|", (progress / (total / 100)));
      pt = millis();
    }
  });
  ArduinoOTA.onError([](ota_error_t error) {
    IPRINTF("Error[%u]: ", error);
    if (error == OTA_AUTH_ERROR) {
      IPRINTLN("Auth Failed");
    } else if (error == OTA_BEGIN_ERROR) {
      IPRINTLN("Begin Failed");
    } else if (error == OTA_CONNECT_ERROR) {
      IPRINTLN("Connect Failed");
    } else if (error == OTA_RECEIVE_ERROR) {
      IPRINTLN("Receive Failed");
    } else if (error == OTA_END_ERROR) {
      IPRINTLN("End Failed");
    }
  });
  ArduinoOTA.begin();
}


#define PORTAL_WAIT_TIMEOUT 30

// wifimanager can run in a blocking mode or a non blocking mode
// Be sure to know how to process loops with no delay() if using non blocking
//bool wm_nonblocking = false; // change to true to use non blocking
WiFiManagerParameter custom_field; // global param ( for non blocking w params )


void setup_wifi(WiFiManager &wm) {
  //void setup_wifi( char* hostname) {
  //pwm = &wm;
  WiFi.mode(WIFI_STA); // explicitly set mode, esp defaults to STA+AP
  //  Serial.begin(115200);
  //  Serial.setDebugOutput(true);
  delay(500);

  pinMode(TRIGGER_PIN, INPUT_PULLUP);
  //  pinMode(LED_BUILTIN, OUTPUT);
  //  digitalWrite(LED_BUILTIN,LOW);

  // wm.resetSettings(); // wipe settings

  //if (wm_nonblocking) wm.setConfigPortalBlocking(false);

  //  // add a custom input field
  //  int customFieldLength = 40;
  //
  //
  //  // new (&custom_field) WiFiManagerParameter("customfieldid", "Custom Field Label", "Custom Field Value", customFieldLength,"placeholder=\"Custom Field Placeholder\"");
  //
  //  // test custom html input type(checkbox)
  //  // new (&custom_field) WiFiManagerParameter("customfieldid", "Custom Field Label", "Custom Field Value", customFieldLength,"placeholder=\"Custom Field Placeholder\" type=\"checkbox\""); // custom html type
  //
  //  // test custom html(radio)
  //  const char* custom_radio_str = "<br/><label for='customfieldid'>Custom Field Label</label><input type='radio' name='customfieldid' value='1' checked> One<br><input type='radio' name='customfieldid' value='2'> Two<br><input type='radio' name='customfieldid' value='3'> Three";
  //  new (&custom_field) WiFiManagerParameter(custom_radio_str); // custom html input
  //
  //  wm.addParameter(&custom_field);
  // wm.setSaveParamsCallback(saveParamCallback);
  //
  //  // custom menu via array or vector
  //  //
  //  // menu tokens, "wifi","wifinoscan","info","param","close","sep","erase","restart","exit" (sep is seperator) (if param is in menu, params will not show up in wifi page!)
  //  // const char* menu[] = {"wifi","info","param","sep","restart","exit"};
  //  // wm.setMenu(menu,6);
  //  std::vector<const char *> menu = {"wifi","info","param","sep","restart","exit"};
  std::vector<const char *> menu = {"wifi", "info", "sep", "restart", "exit"};
  wm.setMenu(menu);

  // set dark theme
  wm.setClass("invert");


  //set static ip
  // wm.setSTAStaticIPConfig(IPAddress(10,0,1,99), IPAddress(10,0,1,1), IPAddress(255,255,255,0)); // set static ip,gw,sn
  // wm.setShowStaticFields(true); // force show static ip fields
  // wm.setShowDnsFields(true);    // force show dns field always

  wm.setConnectTimeout(20); // how long to try to connect for before continuing
  wm.setConfigPortalTimeout(PORTAL_WAIT_TIMEOUT); // auto close configportal after n seconds
  // wm.setCaptivePortalEnable(false); // disable captive portal redirection
  // wm.setAPClientCheck(true); // avoid timeout if client connected to softap

  // wifi scan settings
  // wm.setRemoveDuplicateAPs(false); // do not remove duplicate ap names (true)
  // wm.setMinimumSignalQuality(20);  // set min RSSI (percentage) to show in scans, null = 8%
  // wm.setShowInfoErase(false);      // do not show erase button on info page
  // wm.setScanDispPerc(true);       // show RSSI as percentage not graph icons

  // wm.setBreakAfterConfig(true);   // always exit configportal even if wifi save fails

  bool res;
  // res = wm.autoConnect(); // auto generated AP name from chipid
  // res = wm.autoConnect("AutoConnectAP"); // anonymous ap
  // res = wm.autoConnect("AutoConnectAP","password"); // password protected ap
  IPRINTF("\nAP:%s (%ds) ...\n", device_id_full, PORTAL_WAIT_TIMEOUT);
  res = wm.autoConnect(device_id_full); // ap

  if (!res) {
    IPRINTLN("Failed to connect or hit timeout");
    //ESP.restart();
  }
  else {
    //if you get here you have connected to the WiFi
    IPRINTLN("connected...yeey :)");
  }
}


void checkButton(WiFiManager &wm) {
  //void checkButton(char* hostname) {
  // check for button press
  if ( digitalRead(TRIGGER_PIN) == LOW ) {
    // poor mans debounce/press-hold, code not ideal for production
    delay(50);
    if ( digitalRead(TRIGGER_PIN) == LOW ) {
      IPRINTLN("Button Pressed");
      // still holding button for 3000 ms, reset settings, code not ideaa for production
      delay(3000); // reset delay hold
      if ( digitalRead(TRIGGER_PIN) == LOW ) {
        IPRINTLN("Button Held");
        IPRINTLN("Erasing WIFI Config, restarting");
        //WiFiManager wm;
        wm.resetSettings();
        delay(100);
        ESP.restart();
      }

//      //      // start portal w delay
//      //      IPRINTF("Starting config portal(wait %ds)\n",PORTAL_WAIT_TIMEOUT);
//      IPRINTF("\nAP:%s (%ds) ...\n", device_id_full, PORTAL_WAIT_TIMEOUT);
//      wm.setConfigPortalTimeout(PORTAL_WAIT_TIMEOUT); // auto close configportal after n seconds
//      wm.startConfigPortal(device_id_full);
//      //
//      //      if (!wm.startConfigPortal(hostname)) {
//      //        IPRINTLN("failed to connect or hit timeout");
//      //        delay(3000);
//      //        //ESP.restart();
//      //      } else {
//      //        //if you get here you have connected to the WiFi
//      //        IPRINTLN("connected...yeey :)");
//      //      }
    }
  }
}

//
//String getParam(String name) {
//  WiFiManager &wm = *pwm;
//  //  //read parameter from server, for customhmtl input
//  //  String value;
//  //  if(wm.server->hasArg(name)) {
//  //    value = wm.server->arg(name);
//  //  }
//  //  return value;
//  return "";
//}
//
//void saveParamCallback() {
//  IPRINTLN("[CALLBACK] saveParamCallback fired");
//  //IPRINTLN("PARAM customfieldid = " + getParam("customfieldid"));
//}

//
//void loop() {
//  Serial.printf("loop %s\n",WiFi.status() == WL_CONNECTED?"online":"offline");
//  if(wm_nonblocking == false) delay(1000);
//  if(wm_nonblocking) wm.process(); // avoid delays() in loop when non-blocking and other long running code
//  checkButton();
//  // put your main code here, to run repeatedly:
//}


int publishMQTTMessage(PubSubClient  &client, const char* topic, char* payload, bool retain) {
  led_on();
  bool succeed = false;

  int cut = MQTT_MAX_PACKET_SIZE;                           //拆分字符串发送 //要拆分发送的实际大小

  int payload_len = strlen(payload);                        //总数据长度
  //DPRINTF("payload length = %d\n",payload_len);
  if (payload_len > cut) {
    client.beginPublish(topic, payload_len, retain);        //开始发送长文件参数分别为  主题，长度，是否持续
    int count = payload_len / cut;                          // 2=5/2 2=4/2
    for (int i = 0; i < count; i++) {
      client.write((byte*)(payload + (i * cut)), cut);
    }
    client.write((byte*)(payload + (cut * count)), payload_len - (cut * count));
    succeed = client.endPublish();                          //结束发送文本
  }
  else {
    succeed = client.publish(topic, payload, retain);
  }

  led_off();
  return succeed;
}



void  DEVICE::publish_config(PubSubClient  &client) {
  DPRINTLN();

  // device sensor configs
  if (true) {
    int names_count = 5;

    char names[][9]       = { "Uptime",        "WIFI",       "RSSI",       "IP",       "MAC"      };
    char units[][3]       = { "S",          "",           "dB",         "",           ""           };
    //char dclass[][20]     = { "",           "",           "",           "",         ""           };
    //char sclass[][20]   = {  "",          "",           "",           "",         ""           };
    char categories[][20] = {"diagnostic",  "diagnostic", "diagnostic", "diagnostic", "diagnostic" };



    char read_name[50];
    char unique_id[60];
    char temp[50];

    for (int i = 0; i < names_count; i++) {
      //<discovery_prefix>/<component>/[<node_id>/]<object_id>/config
      sprintf(topic_config, "haworkshopyc1/sensor/%s/%s/config", names[i], device_id);

      sprintf(read_name, "%s", names[i]);
      sprintf(unique_id, "%s_%s_%s", device_manufacture, device_id, names[i]);

      //IPRINTLN(topic_state);
      sprintf(payload, "{\"device\":%s,"
              "\"name\":\"%s\",\"object_id\":\"%s\",\"unique_id\":\"%s\","
              "\"availability_topic\":\"%s\",\"payload_available\":\"online\",\"payload_not_available\":\"offline\","
              "\"state_topic\":\"%s\",\"value_template\":\"{{ value_json.%s }}\","
              , getDeviceConfigJson()
              , read_name, unique_id, unique_id
              , topic_will
              , topic_state_device, names[i]
             );

      //      if (strlen(dclass[i]) != 0)  {
      //        sprintf(temp, "\"device_class\":\"%s\",", dclass[i]);
      //        strcat(payload, temp);
      //      }
      //      if (strlen(sclass[i]) != 0)  {
      //        sprintf(temp, "\"state_class\":\"%s\",", sclass[i]);
      //        strcat(payload, temp);
      //      }
      if (strlen(categories[i]) != 0) {
        sprintf(temp, "\"entity_category\":\"%s\",", categories[i]);
        strcat(payload, temp);
      }
      if (strlen(units[i]) != 0)   {
        sprintf(temp, "\"unit_of_measurement\":\"%s\",", units[i]);
        strcat(payload, temp);
      }

      strcat(payload, "\"retain\":\"false\"}");

      int ret = publishMQTTMessage(client, topic_config, payload, true);
      //      IPRINTF("%s\n", topic_config);
      //      IPRINTF("%s\n", payload);
      DPRINTF("ret(%d) - topic(%2d) - payload(%3d) - %s\n", ret, strlen(topic_config), strlen(payload), topic_config);
    }
    // 数字类
    //number configs
    char _name[] = "Cycle";
    char _unit[] = "S";
    sprintf(topic_config, "haworkshopyc1/number/%s/%s/config",    _name, device_id);
    sprintf(topic_set,    "haworkshopyc1/number/%s/%s/set",     _name, device_id);

    sprintf(read_name, "%s", _name);
    sprintf(unique_id,    "%s_%s_%s", device_manufacture, device_id, _name);

    //\"assumed_state\":\"true\",
    sprintf(payload, "{\"device\":%s,"
            "\"entity_category\":\"config\",\"name\":\"%s\",\"object_id\":\"%s\",\"unique_id\":\"%s\",\"min\":5,\"max\":300,"
            "\"state_topic\":\"%s\",\"value_template\":\"{{ value_json.%s}}\",\"unit_of_measurement \":\"%s\","
            "\"availability_topic\":\"%s\",\"payload_available\":\"online\",\"payload_not_available\":\"offline\","
            "\"command_topic\":\"%s\",\"command_template\":\"{{ value }}\","
            "\"retain\":\"false\"}"
            , getDeviceConfigJson()
            , read_name, unique_id , unique_id
            , topic_state_device , _name, _unit
            , topic_will
            , topic_set);

    int ret = publishMQTTMessage(client, topic_config, payload, true);
    //    IPRINTF("%s\n", topic_config);
    //    IPRINTF("%s\n", payload);
    DPRINTF("ret(%d) - topic(%2d) - payload(%3d) - %s\n", ret, strlen(topic_config), strlen(payload), topic_config);

    // 按钮类
    //button configs
    names_count = 1;
    char _names[][10] = {"Restart"};
    char _dclass[][10] = {"restart"};
    for (int i = 0; i < names_count; i++) {

      sprintf(topic_config, "haworkshopyc1/button/%s/%s/config",    _names[i], device_id);
      sprintf(topic_set,    "haworkshopyc1/button/%s/%s/set",     _names[i], device_id);

      sprintf(read_name, "%s Device", _names[i]);
      sprintf(unique_id,    "%s_%s_%s", device_manufacture, device_id, _names[i]);

      sprintf(payload, "{\"device\":%s,"
              "\"entity_category\":\"config\",\"name\":\"%s\",\"object_id\":\"%s\",\"unique_id\":\"%s\","
              "\"availability_topic\":\"%s\",\"payload_available\":\"online\",\"payload_not_available\":\"offline\","
              "\"command_topic\":\"%s\",\"payload_press\":\"%s\","
              ""
              , getDeviceConfigJson()
              , read_name, unique_id, unique_id
              , topic_will
              , topic_set, _names[i]);

      if (strlen(_dclass[i]) != 0) {
        sprintf(temp, "\"device_class\":\"%s\",", _dclass[i]);
        strcat(payload, temp);
      }
      strcat(payload, "\"retain\":\"false\"}");

      ret = publishMQTTMessage(client, topic_config, payload, true);
      //      IPRINTF("%s\n", topic_config);
      //      IPRINTF("%s\n", payload);
      DPRINTF("ret(%d) - topic(%2d) - payload(%3d) - %s\n", ret, strlen(topic_config), strlen(payload), topic_config);
    }
  }
}



void DEVICE::publish_state(PubSubClient  &client, const  char* wifi_ssid, long rssi) {

  int pos = 0;
  strcpy(payload, "{");
//  for (int i = 0; i < meter.num_cl1xs; i++) {
//    pos = strlen(payload);
//    //DPRINTF("%s link up %d\n",meters[i].name(),meters[i].is_link_up());
//    sprintf(payload + pos, "\"%s\":\"%s\",", meter.cl1xs[i].name(), (meter.cl1xs[i].is_link_up() ? "ok" : "error"));
//  }
  pos = strlen(payload);
  sprintf(payload + pos, "\"Uptime\":%d, \"WIFI\":\"%s\",  \"MAC\":\"%s\",   \"RSSI\":%d,  \"IP\":\"%s\",    \"Cycle\":%d}",
          long(millis() / 1000),   wifi_ssid,       device_mac,        rssi,         device_ip,                 cycle_seconds);

  sprintf(topic_state, "haworkshopyc1/sensor/%s/state", device_id);                                /////////////////////////////////////////TOPIC
  //payload[strlen(payload)] = '\0';

  int ret = publishMQTTMessage(client, topic_state, payload, false);
  if (!ret) IPRINT("Dx");
  else IPRINT("D");

}

void DEVICE::setup(const char* mac) {

  EEPROM_rotate_init();
  byte vv[2];
  EEPROM_rotate_read(vv, 2);
  uint16_t v = vv[0] * 256 + vv[1];
  if (v > 300 || v < 5)    cycle_seconds = 10;
  else cycle_seconds = v;
  IPRINTF("Main loop cycle time: %ds\n", cycle_seconds);

  strcpy(device_mac, mac);
  int j = 0;
  for (int i = 0; i < strlen(device_mac); i++) {
    if (device_mac[i] != ':') {
      device_id[j++] = device_mac[i];
    }
  }
  device_id[j] = '\0';
  sprintf(device_id_full,  "%s-%s-%s", device_manufacture, device_class, device_id);
  sprintf(device_name,     "%s %s %s", device_manufacture, device_class, device_id);
#ifdef ESP32
  sprintf(device_model, "%s %s", device_class,  "3232");
#elif defined(ESP8266)
  sprintf(device_model, "%s %s", device_class,  "8266");
#else
  sprintf(device_model, "%s", device_class);
#endif

  //////////////////////////////////////////////////////
  sprintf(topic_will, "haworkshopyc1/sensor/%s/status", device_id);
  sprintf(topic_state_device, "haworkshopyc1/sensor/%s/state", device_id);

  sprintf(device_config_json, " { \"identifiers\":[\"%s\"],\"manufacturer\":\"%s\",\"sw_version\":\"%s\",\"name\":\"%s\",\"model\":\"%s\",\"connections\":[[\"mac\",\"%s\"]]}"
          , device_id_full, device_manufacture, sw_version, device_name, device_model, device_mac);

}


void DEVICE::setIP(const char* ip) {
  strcpy(device_ip, ip);
}
char* DEVICE::getDeviceID() {
  return device_id_full;
}

char* DEVICE::getDeviceConfigJson() {
  return device_config_json;
}
