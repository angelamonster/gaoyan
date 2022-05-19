
#ifndef DEVICE_H
#define DEVICE_H
#include <WiFiManager.h>
#include <PubSubClient.h>

#define DEBUG   //If you comment this line, the DPRINT & DPRINTLN lines are defined as blank.

#ifdef DEBUG    //Macros are usually in all capital letters.
#define DPRINT(...)    Serial.print(__VA_ARGS__)     //DPRINT is a macro, debug print
#define DPRINTF(...)    Serial.printf(__VA_ARGS__)     //DPRINTF is a macro, debug print
#define DPRINTLN(...)  Serial.println(__VA_ARGS__)   //DPRINTLN is a macro, debug print with new line
#else
#define DPRINT(...)    //now defines a blank line
#define DPRINTF(...)     //now defines a blank line
#define DPRINTLN(...)   //now defines a blank line
#endif


#define IPRINT(...)    Serial.print(__VA_ARGS__)     //DPRINT is a macro, info print
#define IPRINTF(...)    Serial.printf(__VA_ARGS__)     //DPRINTF is a macro, info print
#define IPRINTLN(...)  Serial.println(__VA_ARGS__)   //DPRINTLN is a macro, info print with new line


#ifdef ESP32
#define TRIGGER_PIN 23
#elif defined(ESP8266)
#define TRIGGER_PIN 0
#endif


//
//extern unsigned int pub_signal;
#define PUB_STATE_DEVICE      0b1
#define PUB_STATE_SENSORS     0b10
#define PUB_STATE_SWITCHS     0b100
#define PUB_STATE_ALL         0b1111
#define PUB_CONFIG            0b10000

int publishMQTTMessage(PubSubClient  &client, const char* topic, char* payload, bool retain);
//void mqttcallback(char* intopic, byte* in_payload, unsigned int length);

void setup_OTA();
void setup_wifi(WiFiManager &wm);
void checkButton(WiFiManager &wm);

void led_on();
void led_off();

bool EEPROM_rotate_write(byte* v, int length) ;
bool EEPROM_rotate_read(byte* v, int length);
bool EEPROM_rotate_init();


class DEVICE {
  public:
    void setup(const char* mac);
    void publish_config(PubSubClient  &client);
    void publish_state(PubSubClient  &client,const  char* wifi_ssid,const long rssi);
    
    void mqtt_subscribe_topics(PubSubClient &client);
    bool mqtt_callback(char* intopic, byte* in_payload, unsigned int length);  // return: processed or not

    void setIP(const char* ip);
    char* getDeviceID();
    
    char* getDeviceConfigJson();

    int cycle_seconds;
    unsigned int pub_signal; // SIGNALS FOR UPDATE
};



extern const char sw_version[30];
extern const char device_manufacture[20];               
extern const char device_class[20];

extern char device_id[20];         //    = "XXXXXXXXXXXXXXXXXX";
extern char device_id_full[36];    //  = "KP-PowerMeter-XXXXXXXXXXXXXXXXXXXXXX";
extern char device_name[36];       //  = "KP PowerMeter XXXXXXXXXXXXXXXXXXXXXX";
extern char device_model[36];      //  = "PowerMeter XXCH ESP8266";

extern char device_ip[20];         //  = "000.000.000.000";
extern char device_mac[20];        //  = "XX:XX:XX:XX:XX:XX";

extern char topic_will[50];        //   = "hagz1/sensor/XXXXXXXXXXXXXXXXXXXXXXXXXX/status";
extern char topic_state_device[50];//       = "hagz1/sensor/XXXXXXXXXXXXXXXXXXXXXXXXXX/state";

extern char topic_state[50];       //         = "hagz1/sensor/XXXXXXXXXXXXXXXXXXXX/XXXXXXXXXXXXXXXXXX/state";
extern char topic_config[50];      //        = "hagz1/switch/XXXXXXXXXXXXXXXXXXXX/XXXXXXXXXXXXXXXXXX/config";
extern char topic_set[50];         //           = "hagz1/switch/XXXXXXXXXXXXXXXXXXXX/XXXXXXXXXXXXXXXXXX/set";


extern char payload[1000];
extern char device_config_json[300];


#endif
