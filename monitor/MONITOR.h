#ifndef METER_h
#define METER_h

#include "Arduino.h"
#include <PubSubClient.h>


class MONITOR {
  public:
    void publish_config(PubSubClient  &client);
    void publish_state(PubSubClient  &client);
    void publish_state_switch(PubSubClient  &client);

    void mqtt_subscribe_topics(PubSubClient &client);
    bool mqtt_callback(char* intopic, byte* in_payload, unsigned int length); // return whether is processed

    int workers;
    float unpaid;
    float forcast_24h;

    int ptotal;
    int p[3];
    int t[3]; //temperature
    
    

    bool invalid = true;
};

#endif
