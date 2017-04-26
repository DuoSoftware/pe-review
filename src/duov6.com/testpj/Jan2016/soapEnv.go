package main

import "fmt"
import "duov5.com/DuoAuthorization"

func main() {

	Soap := []byte(`<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
    <s:Body>
        <GetAccessResponse xmlns="http://tempuri.org/">
            <GetAccessResult xmlns:a="http://schemas.datacontract.org/2004/07/DuoAuthSvr" xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
                <a:AccountContracts/>
                <a:AccountRelated>false</a:AccountRelated>
                <a:Application>?</a:Application>
                <a:CompanyID>1</a:CompanyID>
                <a:CompanyIDs xmlns:b="http://schemas.microsoft.com/2003/10/Serialization/Arrays">
                    <b:int>0</b:int>
                    <b:int>101</b:int>
                    <b:int>97</b:int>
                    <b:int>102</b:int>
                    <b:int>105</b:int>
                    <b:int>125</b:int>
                    <b:int>122</b:int>
                    <b:int>130</b:int>
                    <b:int>113</b:int>
                    <b:int>114</b:int>
                    <b:int>103</b:int>
                    <b:int>129</b:int>
                    <b:int>131</b:int>
                    <b:int>132</b:int>
                    <b:int>133</b:int>
                    <b:int>134</b:int>
                    <b:int>135</b:int>
                    <b:int>137</b:int>
                    <b:int>138</b:int>
                    <b:int>139</b:int>
                    <b:int>140</b:int>
                    <b:int>136</b:int>
                    <b:int>141</b:int>
                    <b:int>142</b:int>
                    <b:int>3</b:int>
                    <b:int>9</b:int>
                    <b:int>6</b:int>
                    <b:int>5</b:int>
                    <b:int>17</b:int>
                    <b:int>14</b:int>
                    <b:int>20</b:int>
                    <b:int>21</b:int>
                    <b:int>19</b:int>
                    <b:int>22</b:int>
                    <b:int>18</b:int>
                    <b:int>15</b:int>
                    <b:int>16</b:int>
                    <b:int>90</b:int>
                    <b:int>96</b:int>
                    <b:int>50</b:int>
                    <b:int>94</b:int>
                    <b:int>95</b:int>
                    <b:int>104</b:int>
                    <b:int>100</b:int>
                    <b:int>26</b:int>
                    <b:int>24</b:int>
                    <b:int>28</b:int>
                    <b:int>23</b:int>
                    <b:int>27</b:int>
                    <b:int>25</b:int>
                    <b:int>29</b:int>
                    <b:int>33</b:int>
                    <b:int>30</b:int>
                    <b:int>32</b:int>
                    <b:int>31</b:int>
                    <b:int>35</b:int>
                    <b:int>36</b:int>
                    <b:int>37</b:int>
                    <b:int>39</b:int>
                    <b:int>34</b:int>
                    <b:int>41</b:int>
                    <b:int>48</b:int>
                    <b:int>42</b:int>
                    <b:int>40</b:int>
                    <b:int>99</b:int>
                    <b:int>45</b:int>
                    <b:int>44</b:int>
                    <b:int>46</b:int>
                    <b:int>47</b:int>
                    <b:int>49</b:int>
                    <b:int>52</b:int>
                    <b:int>53</b:int>
                    <b:int>54</b:int>
                    <b:int>55</b:int>
                    <b:int>56</b:int>
                    <b:int>57</b:int>
                    <b:int>58</b:int>
                    <b:int>1</b:int>
                </a:CompanyIDs>
                <a:Data xmlns:b="http://schemas.microsoft.com/2003/10/Serialization/Arrays">
                    <b:KeyValueOfstringstring>
                        <b:Key>Client-IP</b:Key>
                        <b:Value>192.168.5.187</b:Value>
                    </b:KeyValueOfstringstring>
                    <b:KeyValueOfstringstring>
                        <b:Key>Client-Port</b:Key>
                        <b:Value>34111</b:Value>
                    </b:KeyValueOfstringstring>
                </a:Data>
                <a:IgnoreViweObj>false</a:IgnoreViweObj>
                <a:ObjectID>2222</a:ObjectID>
                <a:SecurityToken>5839979680113cda8a66810e691ac63e</a:SecurityToken>
                <a:TenantID>3333</a:TenantID>
                <a:TenantIDs xmlns:b="http://schemas.microsoft.com/2003/10/Serialization/Arrays">
                    <b:int>0</b:int>
                    <b:int>2</b:int>
                    <b:int>15</b:int>
                    <b:int>3</b:int>
                </a:TenantIDs>
                <a:TokenExpireOn>2016-11-16T13:28:20.771</a:TokenExpireOn>
                <a:Type>Admin</a:Type>
                <a:UserName>admin</a:UserName>
                <a:Write>true</a:Write>
                <a:guUserGrpID xmlns:b="http://schemas.microsoft.com/2003/10/Serialization/Arrays">
                    <b:decimal>14102208495058304</b:decimal>
                    <b:decimal>14101511213398405</b:decimal>
                </a:guUserGrpID>
                <a:guUserId>123</a:guUserId>
                <a:viweObjectIDs xmlns:b="http://schemas.microsoft.com/2003/10/Serialization/Arrays">
                    <b:int>1</b:int>
                    <b:int>0</b:int>
                </a:viweObjectIDs>
            </GetAccessResult>
        </GetAccessResponse>
    </s:Body>
</s:Envelope>`)

	res := DuoAuthorization.UserAuth{}

	fmt.Println(res.GetUserAuthObjectFromXML(string(Soap)))

}
