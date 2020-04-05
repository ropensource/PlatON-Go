package wasm.function;

import com.platon.rlp.datatypes.Int32;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ComputeDate;
import org.junit.Before;
import org.junit.Test;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * 验证计算方法在合约里的实现
 * 计算两个日期之间相差的月份
 * @create: 2020/02/16
 * @author liweic
 */

public class ComputeDateTest extends WASMContractPrepareTest {

    private int monthdiff;

    @Before
    public void before() { monthdiff = Integer.parseInt(driverService.param.get("monthdiff"));
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "wasm.ComputeDate验证简单的计算方法",sourcePrefix = "wasm")
    public void Computedate() {

        try {
            prepare();
            ComputeDate computedate = ComputeDate.load("0x2131fcc7630cb13f6e0355df8dc63f2e13191141",web3j, transactionManager, provider);
            String contractAddress = computedate.getContractAddress();

            Int32 datediff1 = computedate.MonthsBetween2Date("20190201", "20200219").send();
            collector.logStepPass("日期月份差:" + datediff1.value);
            collector.assertEqual(datediff1.value, monthdiff);

        } catch (Exception e) {
            collector.logStepFail("ComputeDate failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
