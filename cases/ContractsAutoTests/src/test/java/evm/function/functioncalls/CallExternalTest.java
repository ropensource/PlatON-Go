package evm.function.functioncalls;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.CallExternal;
import org.apache.commons.lang.StringUtils;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;
import java.util.List;


/**
 * @title 验证函数外部调用
 * @description:
 * @author: liweic
 * @create: 2020/01/02 19:11
 **/

public class CallExternalTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "function.CallExternalTest-函数外部调用测试", sourcePrefix = "evm")
    public void callexternal() {
        try {
            CallExternal callexternal = CallExternal.load("0x94fa6b1deafca1fa529e638e26cecf8e398c95f0",web3j, transactionManager, provider);

            String contractAddress = callexternal.getContractAddress();

            //验证函数外部调用
            TransactionReceipt result = callexternal.getResult(new BigInteger("1")).send();
            collector.logStepPass("打印交易Hash：" + result.getTransactionHash());
            collector.logStepPass("intercall函数返回值：" + result);

            List<CallExternal.ExternalCValueEventResponse> eventResult =callexternal.getExternalCValueEvents(result);
            String cv=eventResult.get(0).log.getData();
            collector.assertEqual(subHexData(cv),subHexData("3"));


        } catch (Exception e) {
            collector.logStepFail("CallExternalContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }

    private String subHexData(String hexStr) {
        if (StringUtils.isBlank(hexStr)) {
            throw new IllegalArgumentException("string is blank");
        }
        if (StringUtils.startsWith(hexStr, "0x")) {
            hexStr = StringUtils.substringAfter(hexStr, "0x");
        }
        byte[] addi = hexStr.getBytes();
        for (int i = 0; i < addi.length; i++) {
            if (addi[i] != 0) {
                hexStr = StringUtils.substring(hexStr, i - 1);
                break;
            }
        }
        return hexStr;
    }
}


