
#include "MONITOR.h"
#include "Arduino.h"
#include "HardwareSerial.h"
#include "DEVICE.h"




extern DEVICE device;




void MONITOR::publish_config(PubSubClient  &client ) {
}
//  // meter configs
//  for (int i = 0; i < num_cl1xs; i++) {
//    //meter sensor configs
//
//    int names_count   = 9;
//    char names[][3]      = {"E",                 "F",            "PF",           "U",            "I",            "P",            "Q",              "S",              "T"};
//    char units[][4]       = {"kWh",               "Hz",           "",            "V",            "A",            "W",            "Var",            "VA",             "S"};
//    char dclass[][15]     = {"energy",            "frequency",    "power_factor", "voltage",      "current",      "power",        "reactive_power", "apparent_power", ""};
//    char sclass[][17]     = {"total_increasing",  "measurement",  "measurement",  "measurement",  "measurement",  "measurement",  "measurement",    "measurement",    ""};
//    char categories[][11] = {"",                  "",             "",             "",             "",             "",             "",               "",              ""};
//
//    //sprintf(topic_state, "hagz1/sensor/%s/%s/state", meter.name(), device_id);
//    //sprintf(topic_state, "hagz1/sensor/%s/state", device_id);
//
//    char read_name[50];
//    char switch_name[50];
//    char unique_id[60];
//    char temp[50];
//
//    sprintf(topic_state, "hagz1/sensor/%s/%s/state", cl1xs[i].name(), device_id);
//    for (int j = 0; j < names_count; j++) {
//      //<discovery_prefix>/<component>/[<node_id>/]<object_id>/config
//      sprintf(topic_config, "hagz1/sensor/%s-%s/%s/config", cl1xs[i].name(), names[j], device_id);
//
//      sprintf(read_name, "%s %s", cl1xs[i].name(), names[j]);
//      sprintf(unique_id, "%s_%s_%s_%s", device_manufacture, device_id, cl1xs[i].name(), names[j]);
//
//      //IPRINTLN(topic_state);
//      sprintf(payload, "{\"device\":%s,"
//              "\"name\":\"%s\",\"object_id\":\"%s\",\"unique_id\":\"%s\","
//              //"\"availability_topic\":\"%s\",\"payload_available\":\"online\",\"payload_not_available\":\"offline\","
//              "\"availability_mode\":\"all\",\"availability\":[{\"topic\":\"%s\",\"payload_available\":\"online\",\"payload_not_available\":\"offline\"},{\"topic\":\"%s\",\"value_template\":\"{{ value_json.%s }}\",\"payload_available\":\"ok\",\"payload_not_available\":\"error\"}],"
//              "\"state_topic\":\"%s\",\"value_template\":\"{{ value_json.%s }}\","
//              //"\"json_attributes_topic\":\"%s\",\"json_attributes_template\":\"{{  value_json.%s | tojson}}\","
//              , device.getDeviceConfigJson()
//              , read_name, unique_id, unique_id
//              , topic_will, topic_state_device, cl1xs[i].name()
//              , topic_state,  names[j]
//              //,topic_state, names[j]
//             );
//
//      if (strlen(dclass[j]) != 0)  {
//        sprintf(temp, "\"device_class\":\"%s\",", dclass[j]);
//        strcat(payload, temp);
//      }
//      if (strlen(sclass[j]) != 0)  {
//        sprintf(temp, "\"state_class\":\"%s\",", sclass[j]);
//        strcat(payload, temp);
//      }
//      if (strlen(units[j]) != 0)   {
//        sprintf(temp, "\"unit_of_measurement\":\"%s\",", units[j]);
//        strcat(payload, temp);
//      }
//      if (strlen(categories[j]) != 0) {
//        sprintf(temp, "\"entity_category\":\"%s\",", categories[j]);
//        strcat(payload, temp);
//      }
//
//      strcat(payload, "\"retain\":\"false\",\"expire_after\":360}");
//
//      //DPRINTLN(payload);
//
//      int ret = publishMQTTMessage( client, topic_config, payload, true);
//      //      IPRINTF("%s\n", topic_config);
//      //      IPRINTF("%s\n", payload);
//      //DPRINTF("Sensor topic_config - %s - ret(%d) - topic(%d) - payload(%d)\n", topic_config, ret, strlen(topic_config), strlen(payload));
//      DPRINTF("ret(%d) - topic(%2d) - payload(%3d) - %s\n", ret, strlen(topic_config), strlen(payload), topic_config);
//    }
//
//    //switch configs
//    sprintf(topic_config, "hagz1/switch/%s-ES/%s/config",  cl1xs[i].name(), device_id);
//    sprintf(topic_set,    "hagz1/switch/%s-ES/%s/set",     cl1xs[i].name(), device_id);
//    sprintf(topic_state,  "hagz1/switch/%s/%s/state",   cl1xs[i].name(), device_id);
//
//    sprintf(read_name, "%s %s", cl1xs[i].name(), "Energy State");
//    sprintf(unique_id,    "%s_%s_%s_%s", device_manufacture, device_id, cl1xs[i].name(), "ES");
//
//    sprintf(payload, "{\"device\":%s,"
//            "\"device_class\":\"switch\",\"entity_category\":\"config\",\"name\":\"%s\",\"object_id\":\"%s\",\"unique_id\":\"%s\","
//            "\"state_topic\":\"%s\",\"value_template\":\"{{ value_json.ES }}\","
//            "\"availability_mode\":\"all\",\"availability\":[{\"topic\":\"%s\",\"payload_available\":\"online\",\"payload_not_available\":\"offline\"},{\"topic\":\"%s\",\"value_template\":\"{{ value_json.%s }}\",\"payload_available\":\"ok\",\"payload_not_available\":\"error\"}],"
//            "\"command_topic\":\"%s\",\"payload_on\":\"1\",\"payload_off\":\"0\","
//            "\"assumed_state\":\"true\",\"retain\":\"false\",\"expire_after\":360}"
//            , device.getDeviceConfigJson()
//            , read_name, unique_id, unique_id
//            , topic_state
//            , topic_will, topic_state_device, cl1xs[i].name()
//            , topic_set);
//
//    int ret = publishMQTTMessage( client, topic_config, payload, true);
//    //    IPRINTF("%s\n", topic_config);
//    //    IPRINTF("%s\n", payload);
//    DPRINTF("ret(%d) - topic(%2d) - payload(%3d) - %s\n", ret, strlen(topic_config), strlen(payload), topic_config);
//
//
//    //button configs 按钮
//    sprintf(topic_set,    "hagz1/button/%s-ResetES/%s/set",      cl1xs[i].name(), device_id);
//    sprintf(topic_config, "hagz1/button/%s-ResetES/%s/config",   cl1xs[i].name(), device_id);
//
//    sprintf(read_name, "%s %s", cl1xs[i].name(), "Energy State Reset");
//    sprintf(unique_id,   "%s_%s_%s_%s", device_manufacture, device_id, cl1xs[i].name(), "ResetES");
//
//    sprintf(payload, "{\"device\":%s,"
//            "\"entity_category\":\"diagnostic\",\"device_class\":\"update\",\"name\":\"%s\",\"object_id\":\"%s\",\"unique_id\":\"%s\","
//            "\"availability_mode\":\"all\",\"availability\":[{\"topic\":\"%s\",\"payload_available\":\"online\",\"payload_not_available\":\"offline\"},{\"topic\":\"%s\",\"value_template\":\"{{ value_json.%s }}\",\"payload_available\":\"ok\",\"payload_not_available\":\"error\"}],"
//            "\"command_topic\":\"%s\",\"payload_press\":\"reset\","
//            "\"retain\":\"false\"}"
//            , device.getDeviceConfigJson()
//            , read_name, unique_id, unique_id
//            , topic_will, topic_state_device, cl1xs[i].name()
//            , topic_set);
//
//    ret = publishMQTTMessage( client, topic_config, payload, true);
//    //    IPRINTF("%s\n", topic_config);
//    //    IPRINTF("%s\n", payload);
//    DPRINTF("ret(%d) - topic(%2d) - payload(%3d) - %s\n", ret, strlen(topic_config), strlen(payload), topic_config);
//  }
//  IPRINT("C");
//}

void MONITOR::publish_state(PubSubClient  &client ) {

  //  for (int i = 0; i < num_cl1xs; i++) {
  //
  //    float E     = 0.0;
  //    float U     = 0.0;
  //    float I     = 0.0;
  //    float P     = 0.0;
  //    float F     = 0.0;
  //    float S     = 0.0;
  //
  //    bool succeed = true;
  //
  //    succeed = succeed & cl1xs[i].getE(&E);
  //    E = E / 1000.0;                 // Wh to KWH
  //    succeed = succeed & cl1xs[i].getU(&U);
  //    succeed = succeed & cl1xs[i].getI(&I);
  //    succeed = succeed & cl1xs[i].getP(&P);
  //    succeed = succeed & cl1xs[i].getF(&F);
  //    succeed = succeed & cl1xs[i].getS(&S);
  //
  //    float Q     =  0.0;
  //    float PF     =  0.0;
  //    float T     =  0.0;
  //    succeed = succeed & cl1xs[i].getQ(&Q);
  //    succeed = succeed & cl1xs[i].getPF(&PF);
  //    succeed = succeed & cl1xs[i].getT(&T);
  //
  //    int ES     = 0;
  //    succeed = succeed & cl1xs[i].getEnergyState(&ES);
  //    //#ifdef DEBUG
  //    //    char name[10];
  //    //    char ver[10];
  //    //    succeed = succeed & cl1xs[i].getModuleName(name);
  //    //    succeed = succeed & cl1xs[i].getVersion(ver);
  //    //#endif
  //    //    DPRINTF("%s - %s,%s,U:%.2f,I:%.2f,P:%.2f,PF:%.2f,F:%.2f,Q:%.2f,S:%.2f,E:%.2f,T:%.2f\n", cl1xs[i].name(), name, ver, U, I, P, PF, F, Q, S, E, T);
  //
  //
  //    if (succeed) {
  //      sprintf(payload, "{\"E\":%.2f,\"F\":%.2f,\"PF\":%.2f,\"U\":%.2f,\"I\":%.2f,\"P\":%.2f,\"Q\":%.2f,\"S\":%.2f,\"T\":%.0f}",
  //              E,           F,          PF,         U,         I,         P,         Q,         S,       T);
  //
  //      sprintf(topic_state, "hagz1/sensor/%s/%s/state", cl1xs[i].name(), device_id);
  //
  //      int ret = publishMQTTMessage( client, topic_state, payload, false);
  //      if (!ret) IPRINTF("Failed to pulish - %s\n", topic_state);
  //      else IPRINT("S");
  //
  //    } else       IPRINTF("s");
  //
  //  }

}

void MONITOR::publish_state_switch(PubSubClient  &client ) {
  //  for (int i = 0; i < num_cl1xs; i++) {
  //
  //    int ES     = 0;
  //    bool succeed = cl1xs[i].getEnergyState(&ES);
  //
  //    if (succeed) {
  //      sprintf(payload, "{\"ES\":%d}", ES);
  //      sprintf(topic_state, "hagz1/switch/%s/%s/state", cl1xs[i].name(), device_id);
  //
  //      int ret = publishMQTTMessage( client, topic_state, payload, false);
  //      if (!ret) IPRINTF("Failed to pulish - %s\n", topic_state);
  //      else IPRINT("W");
  //    } else  IPRINTF("w");
  //  }
}




void MONITOR::mqtt_subscribe_topics(PubSubClient &client) {
  int succeed;

  sprintf(topic_set,    "haworkshopyc1/sensor/powermeteryc1/state");
  succeed = client.subscribe(topic_set);
  IPRINTF("   - %d - %s\n",succeed, topic_set);

  char* hs[10] = {"w0004", "w0005", "w0007"};
  for (int i = 0; i < 3; i++) {
    sprintf(topic_set,    "haworkshopyc1/sensor/%s/state", hs[i]);
    client.subscribe(topic_set);
    IPRINTF("   - %d - %s\n",succeed, topic_set);
  }

  
  sprintf(topic_set,    "haworkshopyc1/sensor/hive/state");
  client.subscribe(topic_set);
  IPRINTF("   - %d - %s\n",succeed, topic_set);
}


#include <ArduinoJson.h>                              /* 引入JSON解析需要用的库文件 */
DynamicJsonDocument doc(1024);


bool MONITOR::mqtt_callback(char* intopic, byte* in_payload, unsigned int length) {  // return whether is processed
  strncpy(payload, (char*)in_payload, length);
  payload[length] = '\0';

  sprintf(topic_set,    "haworkshopyc1/sensor/hive/state");
  if (!strcmp(intopic, topic_set)) {
    deserializeJson(doc, payload);
    const char* count_str = doc["w"];
    workers = atoi(count_str);
    //const char* unpaid_str = doc["totalUnpaid"];
    //unpaid = atof(unpaid_str);
    unpaid = doc["tup"];
    forcast_24h = doc["e24h"];
    invalid = true;
    IPRINT("P");
    return true;
  }

  sprintf(topic_set,    "haworkshopyc1/sensor/powermeteryc1/state");
  if (!strcmp(intopic, topic_set)) {
    deserializeJson(doc, payload);
//    p[0] = doc["p"][0];
//    p[1] = doc["p"][1];
//    p[2] = doc["p"][2];
    ptotal = doc["P"];
    invalid = true;
    IPRINT("M");
    return true;
  }

  char* hs[10] = {"w0004", "w0005", "w0007"};
  for (int i = 0; i < 3; i++) {
    sprintf(topic_set,    "haworkshopyc1/sensor/%s/state", hs[i]);
    if (!strcmp(intopic, topic_set)) {
      deserializeJson(doc, payload);
      t[i] = doc["ht"];
      invalid = true;
      //DPRINTLN(payload);
    IPRINTF("%d",i);
      return true;
    }

  }

  return false;
}
