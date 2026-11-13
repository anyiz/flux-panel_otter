package com.admin;

import com.admin.common.utils.GostUtil;
import com.admin.entity.Tunnel;
import com.alibaba.fastjson2.JSON;
import com.alibaba.fastjson2.JSONObject;
import org.junit.jupiter.api.Test;
import java.lang.reflect.Method;
import static org.junit.jupiter.api.Assertions.*;

class GostProxyProtocolWireTest {
    private JSONObject wire(String protocol, Integer version, int type) throws Exception {
        Tunnel tunnel = new Tunnel();
        tunnel.setTcpListenAddr("0.0.0.0");
        tunnel.setUdpListenAddr("0.0.0.0");
        Method method = GostUtil.class.getDeclaredMethod("createServiceConfig",
            String.class, Integer.class, Integer.class, String.class, String.class,
            Integer.class, Tunnel.class, String.class, String.class, Integer.class);
        method.setAccessible(true);
        Object config = method.invoke(null, "test", 65421, null, "127.0.0.1:18001",
            protocol, type, tunnel, "fifo", null, version);
        return JSON.parseObject(JSON.toJSONString(config));
    }

    @Test
    void proxyVersionSurvivesJsonAsAStringForGost() throws Exception {
        for (int type : new int[]{1, 2}) {
            for (int version : new int[]{1, 2}) {
                Object value = wire("tcp", version, type).getJSONObject("handler")
                    .getJSONObject("metadata").get("proxyProtocol");
                assertEquals(String.valueOf(version), value);
            }
        }
    }

    @Test
    void disabledAndUdpRulesDoNotSendProxyHeaders() throws Exception {
        for (Integer version : new Integer[]{null, 0}) {
            assertNull(wire("tcp", version, 1).getJSONObject("handler").get("metadata"));
        }
        assertNull(wire("udp", 2, 1).getJSONObject("handler").get("metadata"));
    }
}
